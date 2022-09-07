package prompts

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"tie.prodigy9.co/config"

	"golang.org/x/crypto/ssh/terminal"
)

type Session struct {
	cfg     *config.Config
	scanner *bufio.Scanner

	args []string
}

func New(ctx *config.Config, args []string) *Session {
	return &Session{ctx, nil, args}
}

func (s *Session) Len() int {
	return len(s.args)
}

func (s *Session) YesNo(question string) bool {
	if s.cfg.AlwaysYes() {
		return true
	}

	answer := s.readYesNo(question)
	answer = strings.TrimSpace(answer)
	answer = strings.ToUpper(answer)
	switch answer {
	case "1", "Y", "YES":
		return true
	default:
		return false
	}
}

func (s *Session) SensitiveStr(item string) string {
	if len(s.args) > 0 {
		head, tail := s.args[0], s.args[1:]
		s.args = tail
		return head

	} else {
		input := s.readSensitiveInput(item)
		fmt.Fprintln(os.Stdout)
		return input

	}
}

func (s *Session) OptionalStr(item string, defaultValue string) string {
	if len(s.args) > 0 {
		return s.Str(item)
	} else {
		return defaultValue
	}
}

func (s *Session) Str(item string) string {
	if len(s.args) > 0 {
		head, tail := s.args[0], s.args[1:]
		s.args = tail
		return head

	} else {
		return s.readInput(item)

	}
}

func (s *Session) readYesNo(question string) string {
	if s.scanner == nil {
		s.scanner = bufio.NewScanner(os.Stdin)
	}

	fmt.Fprintf(os.Stderr, question+" (y/n)? ")
	return s.mustScan()
}

func (s *Session) readSensitiveInput(item string) string {
	fmt.Fprintf(os.Stderr, "enter "+item+" securely: ")
	return s.mustScanSensitive()
}

func (s *Session) readInput(item string) string {
	if s.scanner == nil {
		s.scanner = bufio.NewScanner(os.Stdin)
	}

	fmt.Fprintf(os.Stderr, "enter "+item+": ")
	return s.mustScan()
}

func (s *Session) mustScanSensitive() string {
	bytes, err := terminal.ReadPassword(0)
	if err != nil {
		s.cfg.Fatalln("i/o error", err)
		return ""
	}

	return string(bytes)
}

func (s *Session) mustScan() string {
	if !s.scanner.Scan() {
		if s.scanner.Err() != nil {
			s.cfg.Fatalln("i/o error", s.scanner.Err())
			return ""
		} else {
			s.cfg.Fatalln("expect more input on stdin")
			return ""
		}
	}

	return s.scanner.Text()
}
