package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/sevlyar/go-daemon"
	"github.com/urfave/cli/v2"
	"golang.org/x/term"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

var (
	username   string
	password   string
	gap        int64
	loginUrl   string
	testingUrl string
	isLoop     bool
	isDaemon   bool
)

const (
	UsernameEnvName string = "GO_DRCOM_USERNAME"
	PasswordEnvName string = "GO_DRCOM_PASSWORD"
)

func GbkToUtf8(str []byte) (b []byte, err error) {
	r := transform.NewReader(bytes.NewBuffer(str), simplifiedchinese.GBK.NewDecoder())
	b, err = ioutil.ReadAll(r)
	return
}

func Login(username, password string) error {
	log.Print("try login now")
	payload := fmt.Sprintf("DDDDD=%s&upass=%s&R1=0&R2=&R6=0&para=00&0MKKey=123456", url.QueryEscape(username), url.QueryEscape(password))
	// fmt.Println(payload)
	resp, err := http.Post(loginUrl, "application/x-www-form-urlencoded", bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	utf8Body, err := GbkToUtf8(body)
	if err != nil {
		return err
	}

	// fmt.Println(string(utf8Body))

	if !strings.Contains(string(utf8Body), "<title>登录成功窗</title>") && !strings.Contains(string(utf8Body), "<title>Drcom PC登陆成功页</title>") {
		return fmt.Errorf("login fail")
	}

	return nil
}

func IsLogin(url string) bool {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resq, err := client.Get(url)
	if err != nil {
		log.Printf("login testing fail and return %#v", err)
		return false
	}

	// defer resq.Body.Close()
	// body, err := ioutil.ReadAll(resq.Body)
	// if err != nil {
	// 	return false
	// }
	// fmt.Println(string(body))

	if resq.StatusCode != 200 {
		log.Printf("login testing fail and status code is %#v", resq.StatusCode)
		return false
	}

	log.Print("login testing is success")
	return true
}

func LoginLoop(username, password, testingUrl string, gap int64) {

	for {
		if !IsLogin(testingUrl) {
			err := Login(username, password)
			if err != nil {
				log.Printf("Login fail and error is %#v", err)
			}
		}
		time.Sleep(time.Duration(gap) * time.Second)
	}
}

func Contains(tokens []string, token string) bool {
	for _, t := range tokens {
		if token == t {
			return true
		}
	}
	return false
}

func BecomdDaemon() *daemon.Context {
	log.Printf("become daemon")
	if runtime.GOOS == "windows" {
		fmt.Errorf("You are running in windows.\nDaemon is not support now.\n")
		return nil
	}
	args := os.Args
	// log.Println(args)
	envs := os.Environ()
	envs = append(envs, fmt.Sprintf("%s=%s", UsernameEnvName, username))
	envs = append(envs, fmt.Sprintf("%s=%s", PasswordEnvName, password))
	cntxt := &daemon.Context{
		PidFileName: "go-drcom.pid",
		PidFilePerm: 0644,
		LogFileName: "go-drcom.log",
		LogFilePerm: 0640,
		WorkDir:     "./",
		Umask:       027,
		Args:        args,
		Env:         envs,
	}

	d, err := cntxt.Reborn()
	if err != nil {
		log.Fatal("Unable to run: ", err)
	}
	if d != nil {
		return nil
	}

	log.Print("- - - - - - - - - - - - - - -")
	log.Print("daemon started")
	return cntxt
}

func ActionFunc(c *cli.Context) error {
	if username == "" {
		fmt.Print("Username: ")
		fmt.Scanln(&username)
	}
	if password == "" {
		fmt.Print("Password: ")
		passwordByte, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			log.Fatalf("read password fail, %#v", err.Error())
		}
		password = string(passwordByte)
		fmt.Println()
	}

	if isDaemon && isLoop {
		cntxt := BecomdDaemon()
		if cntxt != nil {
			defer cntxt.Release()
		} else {
			return nil
		}
	} else if isDaemon && !isLoop {
		fmt.Errorf("Warning: isLoop is false, isDaemon is disable now...")
	}

	if isLoop {
		LoginLoop(username, password, testingUrl, gap)
		return nil
	}
	return Login(username, password)

}

func main() {

	app := &cli.App{
		Name:  "go-drcom",
		Usage: "A simple and easy easy-to-use command line app to login drcom (in szu).",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "username",
				Aliases:     []string{"u"},
				Usage:       "Username for drcom",
				Destination: &username,
				Value:       "",
				EnvVars:     []string{UsernameEnvName},
				// Required:    true,
			},
			&cli.StringFlag{
				Name:        "password",
				Aliases:     []string{"p"},
				Usage:       "Password for drcom",
				Destination: &password,
				Value:       "",
				EnvVars:     []string{PasswordEnvName},
				// Required:    true,
			},
			&cli.Int64Flag{
				Name:        "gap",
				Aliases:     []string{"g"},
				Usage:       "Time gap for check network",
				Destination: &gap,
				Value:       60,
			},
			&cli.StringFlag{
				Name:        "loginUrl",
				Aliases:     []string{"l"},
				Usage:       "Login Url",
				Destination: &loginUrl,
				Value:       "https://drcom.szu.edu.cn",
			},
			&cli.StringFlag{
				Name:        "testingUrl",
				Aliases:     []string{"t"},
				Usage:       "Testing Url",
				Destination: &testingUrl,
				Value:       "https://www.baidu.com",
			},
			&cli.BoolFlag{
				Name:        "isLoop",
				Aliases:     []string{"loop"},
				Usage:       "Login Loop",
				Destination: &isLoop,
				Value:       false,
			},
			&cli.BoolFlag{
				Name:        "isDaemon",
				Aliases:     []string{"daemon", "d"},
				Usage:       "Whether daemon",
				Destination: &isDaemon,
				Value:       false,
			},
		},
		Action: ActionFunc,
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

	err := app.Run(os.Args)
	if err != nil {
		log.Print(err)
	}
}
