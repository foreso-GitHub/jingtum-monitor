package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
)

type Request struct {
	Id      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  Params `json:"params"`
}

type Params struct {
	User     string `json:"user"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Command  string `json:"command"`
	Mode     int    `json:"mode"`
}

func main() {
	fmt.Println("===start===")
	//http.HandleFunc("/jsonrpc", jsonrpc)
	//////http.ListenAndServe("0.0.0.0:9999",nil)
	//http.ListenAndServe("127.0.0.1:9545",nil)

	params := LoadParams("./params.json")
	if params.Mode == 1 {
		runSsh(params)
	} else {
		runSshT(params)
	}
	fmt.Println(params.Command + " done!")

	//params := new(Params)
	//
	////params.User = "root"
	////params.Password = ""
	////params.Host = "39.99.174.194"
	////params.Port = 22
	////params.Command = "/work/test2.sh"
	////runSsh(*params)
	////fmt.Println(params.Host + " done!")
	////
	//params.User = "jt"
	//params.Password = "jt"
	////params.Host = "box-admin.elerp.net"
	////params.Port = 2003
	//params.Host = "10.0.0.203"
	//params.Port = 22
	//params.Command = "/home/jt/test.sh"
	//runSsh(*params)
	//fmt.Println(params.Command + " done!")
	//
	//params = new(Params)
	//params.User = "jt"
	//params.Password = "jt"
	////params.Host = "box-admin.elerp.net"
	////params.Port = 2003
	//params.Host = "10.0.0.203"
	//params.Port = 22
	//params.Command = "/home/jt/start.sh"
	//runSsh(*params)
	//fmt.Println(params.Command + " done!")

	//params.Command = "/home/jt/stop.sh"
	//runSsh(*params)
	//fmt.Println(params.Command + " done!")
	//
	//params.Command = "/home/jt/clear.sh"
	//runSsh(*params)
	//fmt.Println(params.Command + " done!")

	//conclusion:
	//1. can exec test.sh of node 203, from node 202.
	//2. cannot invoke api from port 9545
	//3. mode1 can work if user ignore the password.  to do it, vi /etc/sudoers, change to "jt      ALL=(ALL) NOPASSWD: ALL".

	fmt.Println("===end===")

}

func LoadParams(path string) Params {
	var params Params
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panicln("load config file failed: ", err)
	}
	err = json.Unmarshal(buf, &params)
	if err != nil {
		log.Panicln("decode config file failed:", string(buf), err)
	}
	return params
}

func runCmd(user string, pw string, host string, port int, cmd string) {
	params := new(Params)
	params.User = user
	params.Password = pw
	params.Host = host
	params.Port = port
	params.Command = cmd
	//runSsh(*params)
	runSshT(*params)
	fmt.Println(params.Command + " done!")
}

func jsonrpc(res http.ResponseWriter, req *http.Request) {
	//读取body
	body, _ := ioutil.ReadAll(req.Body)
	strBody := string(body)
	fmt.Println("bodyStr", strBody)
	requestObject := new(Request)
	//json.Unmarshal(body, &requestObject)
	err := json.Unmarshal(body, &requestObject)
	if err != nil {
		log.Panicln("decode request failed:", string(body), err)
	}

	fmt.Println("requestObject: ", *requestObject)
	//fmt.Fprint(res, requestObject.Method)

	success, err := dealMethod(requestObject)
	result := "success!"
	if !success {
		result = err.Error()
	}

	//result = runSsh("/work/test2.sh")
	result = runSsh(requestObject.Params)

	fmt.Fprint(res, result)
}

func dealMethod(request *Request) (bool, error) {
	re, err := Cmd("echo hello word")
	fmt.Println(re, err)
	//re,err = CmdTime("bash /root/workplace/src/go-exec/i.sh",5)
	//fmt.Println(re,err)
	if err != nil {
		return false, err
	}
	return true, nil
}

func Cmd(command string) (string, error) {
	cmd := exec.Command("bash", "-c", command)
	out, err := cmd.CombinedOutput()
	result := string(out)
	return result, err
}

func SSHConnect(user, password, host string, port int) (*ssh.Session, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		client       *ssh.Client
		session      *ssh.Session
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	hostKeyCallbk := func(hostname string, remote net.Addr, key ssh.PublicKey) error {
		return nil
	}

	clientConfig = &ssh.ClientConfig{
		User: user,
		Auth: auth,
		// Timeout:             30 * time.Second,
		HostKeyCallback: hostKeyCallbk,
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if client, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create session
	if session, err = client.NewSession(); err != nil {
		return nil, err
	}

	return session, nil
}

func runSsh(params Params) string {

	var stdOut, stdErr bytes.Buffer

	session, err := SSHConnect(params.User, params.Password, params.Host, params.Port)
	//session, err := SSHConnect( "jt", "jt", "180.152.251.68", 10202 )
	//session, err := SSHConnect( "jt", "jt", "box-admin.elerp.net", 10202 )
	//session, err := SSHConnect( "root", "", "39.99.174.194", 22 )
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	session.Stdout = &stdOut
	session.Stderr = &stdErr
	//session.Stdout = os.Stdout
	//session.Stderr = os.Stderr

	//session.Run("if [ -d /root ]; then echo 0; else echo 1; fi")
	//session.Run("if [ pwd ]; then echo 0; else echo 1; fi")
	//session.Run("#!/bin/bash;	echo \"Hello World !\"")
	session.Run(params.Command)
	//result, err := session.Output("ls -al")
	//ret, err := strconv.Atoi( strings.Replace( stdOut.String(), "\n", "", -1 ))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%d, %s\n", ret, stdErr.String())
	fmt.Println(stdOut.String())
	return stdOut.String()

	//fmt.Println(os.Stdout)
	//return "os.Stdout"
}

func runSshT(params Params) string {

	//var stdOut, stdErr bytes.Buffer

	//session, err := SSHConnect( params.User, params.Password, params.Host, params.Port )
	////session, err := SSHConnect( "jt", "jt", "180.152.251.68", 10202 )
	////session, err := SSHConnect( "jt", "jt", "box-admin.elerp.net", 10202 )
	////session, err := SSHConnect( "root", "", "39.99.174.194", 22 )
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer session.Close()

	session, err := SSHConnect(params.User, params.Password, params.Host, params.Port)
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	fd := int(os.Stdin.Fd())
	oldState, err := terminal.MakeRaw(fd)
	if err != nil {
		panic(err)
	}
	defer terminal.Restore(fd, oldState)

	//session.Stdout = &stdOut
	//session.Stderr = &stdErr
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		panic(err)
	}

	// Set up terminal modes
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,     // enable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	// Request pseudo terminal
	if err := session.RequestPty("xterm-256color", termHeight, termWidth, modes); err != nil {
		log.Fatal(err)
	}

	//session.Run("if [ -d /root ]; then echo 0; else echo 1; fi")
	//session.Run("if [ pwd ]; then echo 0; else echo 1; fi")
	//session.Run("#!/bin/bash;	echo \"Hello World !\"")
	session.Run(params.Command)
	//result, err := session.Output("ls -al")
	//ret, err := strconv.Atoi( strings.Replace( stdOut.String(), "\n", "", -1 ))
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Printf("%d, %s\n", ret, stdErr.String())
	//fmt.Println(stdOut.String())
	//return stdOut.String()

	fmt.Println(os.Stdout)
	return "os.Stdout"
}
