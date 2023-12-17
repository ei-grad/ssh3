package unix_util

import (
	"io"
	"os/exec"
	osuser "os/user"
	"strconv"
	"syscall"

	"github.com/rs/zerolog/log"
)


type User struct {
	Username string
	Uid		 uint64
	Gid		 uint64
	Dir		 string
	Shell	 string
}

func GetUser(username string) (*User, error) {
	u, err := osuser.Lookup(username)
	if err != nil {
		return nil, err
	}

	uid, err := strconv.Atoi(u.Uid)
	if err != nil {
		log.Error().Msgf("could not convert uid %s into an int", u.Uid)
		return nil, err
	}
	gid, err := strconv.Atoi(u.Uid)
	if err != nil {
		log.Error().Msgf("could not convert gid %s into an int", u.Uid)
		return nil, err
	}

	return &User{
		Username: u.Username,
		Uid: uint64(uid),
		Gid: uint64(gid),
		Dir: u.HomeDir,
		Shell: "/bin/bash",
	}, nil
}

func (u *User) CreateShell(addEnv string, stdout, stderr io.Writer, stdin io.Reader) (*exec.Cmd, error) {
	cmd, _, _, _, err := u.CreateCommand(addEnv, stdout, stderr, stdin, u.Shell)
	return cmd, err
}


func (u *User) CreateCommand(addEnv string, stdout, stderr io.Writer, stdin io.Reader, command string, args ...string) (*exec.Cmd, io.Reader, io.Reader, io.Writer, error) {
	cmd := exec.Command(command, args...)
	
	cmd.Env = append(cmd.Env, addEnv)
	cmd.Dir = u.Dir

	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(u.Uid), Gid: uint32(u.Gid)}

	var err error
	var stdoutR, stderrR io.Reader
	var stdinW io.Writer

	if stdout == nil {
		stdoutR, err = cmd.StdoutPipe()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	} else {
		cmd.Stdout = stdout
	}
	if stderr == nil {
		stderrR, err = cmd.StderrPipe()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	} else {
		cmd.Stderr = stderr
	}
	if stdin == nil {
		stdinW, err = cmd.StdinPipe()
		if err != nil {
			return nil, nil, nil, nil, err
		}
	} else {
		cmd.Stdin = stdin
	}

	return cmd, stdoutR, stderrR, stdinW, err
}

func (u *User) CreateCommandPipeOutput(addEnv string, command string, args ...string) (*exec.Cmd, io.Reader, io.Reader, io.Writer, error) {
	cmd := exec.Command(command, args...)
	
	cmd.Env = append(cmd.Env, addEnv)
	cmd.Dir = u.Dir

	cmd.SysProcAttr = &syscall.SysProcAttr{}
	cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(u.Uid), Gid: uint32(u.Gid)}

	return u.CreateCommand(addEnv, nil, nil, nil, command, args...)
}


/*
 *  Returns a boolean stating whether the user is correctly authenticated on this
 *  server. May return a UserNotFound error when the user does not exist.
 */
 func UserPasswordAuthentication(username, password string) (bool, error) {
	return userPasswordAuthentication(username, password)
}