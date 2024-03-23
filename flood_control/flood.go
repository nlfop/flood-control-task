package flood

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type InfoUser struct {
	Num       int
	TimeStart int64
}

type FloodManager struct {
	DB *sql.DB
	N  int64
	K  int
}

var (
	ErrNoAuth = errors.New("no user found")
	ErrN      = errors.New("n is missing")
	ErrK      = errors.New("k is missing")
)

func InitDBFlood(ctx context.Context) (*FloodManager, error) {
	fm := &FloodManager{}
	err := fm.ChekNK()
	if err != nil {
		return nil, err
	}
	dsn := "root:12345Anast@tcp(localhost:3306)/golang?"
	dsn += "charset=utf8"
	dsn += "&interpolateParams=true"

	fm.DB, err = sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	fm.DB.SetMaxOpenConns(10)
	err = fm.DB.Ping()
	if err != nil {
		return nil, err
	}
	go fm.DeleteOldSession(ctx)

	return fm, nil
}

func (fm *FloodManager) Check(ctx context.Context, userID int64) (bool, error) {
	timeNow := time.Now().Unix()

	err := fm.ChekNK()
	if err != nil {
		return false, err
	}

	infoU, err := fm.getUserInfo(userID)
	fmt.Println(infoU)
	if err == ErrNoAuth {

		if (timeNow - infoU.TimeStart) > fm.N {
			infoU = &InfoUser{
				Num:       1,
				TimeStart: timeNow,
			}
			fm.putUserInfo(infoU, userID)
			return false, nil
		}

		if infoU.Num >= fm.K {
			return true, nil
		}

		infoU.Num++

	} else if err == nil {

		infoU = &InfoUser{
			Num:       1,
			TimeStart: timeNow,
		}
	}

	fm.putUserInfo(infoU, userID)
	return false, err
}

func (fm *FloodManager) getUserInfo(userID int64) (*InfoUser, error) {
	infoU := &InfoUser{}
	err := fm.DB.
		QueryRow("SELECT numr, timestart FROM flood WHERE id = ?", userID).
		Scan(&infoU.Num, &infoU.TimeStart)
	if err != nil {
		return &InfoUser{}, ErrNoAuth
	}

	return infoU, nil
}

func (fm *FloodManager) putUserInfo(infoU InfoUser, userID int64) error {
	_, err := fm.DB.Exec(
		"INSERT INTO flood (`id`, `numr`, `timestart`) VALUES (?, ?, ?)",
		userID,
		infoU.Num,
		infoU.TimeStart,
	)
	if err != nil {

		_, err = fm.DB.Exec(
			"UPDATE flood SET `numr` = ?, `timestart` = ?  WHERE id = ?",
			infoU.Num, infoU.TimeStart, userID,
		)
		if err != nil {
			return err
		}
	}

	return err
}

func (fm *FloodManager) DeleteOldSession(ctx context.Context) {

	for {
		timer1 := time.NewTimer(1 * time.Minute)

		select {
		case <-timer1.C:
			rows, err := fm.DB.QueryContext(ctx, "SELECT id, timestart FROM floods")
			if err != nil {
				panic(err)
			}
			for rows.Next() {
				var idUser, timeDB int64
				err = rows.Scan(&idUser, &timeDB)
				if err != nil {
					panic(err)
				}
				if (time.Now().Unix() - timeDB) > fm.N {
					_, err := fm.DB.Exec(
						"DELETE FROM floods WHERE id = ?",
						idUser,
					)
					if err != nil {
						fmt.Println(err)
						return
					}

				}
			}
			rows.Close()

		case <-ctx.Done():
			timer1.Stop()
			return
		}

	}
}

func (fm *FloodManager) ChekNK() error {
	var err error
	numInit, exists := os.LookupEnv("NFLOOD")
	if !exists {
		return ErrN
	}
	keyInit, exists := os.LookupEnv("KFLOOD")
	if !exists {
		return ErrK
	}
	fm.N, err = strconv.ParseInt(numInit, 10, 64)
	if err != nil {
		return err
	}
	fm.K, err = strconv.Atoi(keyInit)
	if err != nil {
		return err
	}
	return nil
}
