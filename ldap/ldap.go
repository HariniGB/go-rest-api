package ldap

import (
	"fmt"
	"log"
	"strings"

	"github.com/HariniGB/login-provider/models"
	"gopkg.in/ldap.v2"
)

type Ldap struct {
	username string
	password string
	host     string
	port     int
	dn       string
}

func NewLdap(username, password, host string, port int, dn string) (*Ldap, error) {
	l := &Ldap{
		username: username,
		password: password,
		host:     host,
		port:     port,
		dn:       dn,
	}

	if err := l.Test(); err != nil {
		return nil, err
	}

	l.createTLD("groups")
	l.createTLD("people")

	return l, nil
}

func (l *Ldap) Test() error {
	c, err := l.newConnection()
	if err != nil {
		return err
	}
	defer c.Close()
	return nil
}

func (l *Ldap) newConnection() (*ldap.Conn, error) {
	c, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.host, l.port))
	if err != nil {
		return nil, err
	}

	return c, nil
}

func (l *Ldap) ExistsUser(name string) bool {
	conn, err := l.newConnection()
	if err != nil {
		return false
	}
	defer conn.Close()
	dn := fmt.Sprintf("cn=%s,ou=people,%s", name, l.dn)
	_, err = l.search(conn, dn, fmt.Sprintf("(&(objectclass=person)(cn=%s))", name), []string{"cn"})
	if err != nil {
		return false
	}

	return true
}

func (l *Ldap) AddUser(user *models.User) error {
	conn, err := l.loginAdmin(l.username, l.password)
	if err != nil {
		return err
	}
	defer conn.Close()

	dn := fmt.Sprintf("cn=%s,ou=people,%s", user.Username, l.dn)
	req := ldap.NewAddRequest(dn)
	req.Attribute("description", []string{fmt.Sprintf("%s user", user.Username)})
	req.Attribute("userPassword", []string{user.Password})
	req.Attribute("sn", []string{user.LastName})
	req.Attribute("objectclass", []string{"iNetOrgPerson"})
	req.Attribute("givenname", []string{user.FirstName})
	req.Attribute("mail", []string{user.Email})
	err = conn.Add(req)

	if ldap.IsErrorWithCode(err, 68) {
		return UserAlreadyExists
	}
	return err
}

func (l *Ldap) DeleteUser(name string) error {
	conn, err := l.loginAdmin(l.username, l.password)
	if err != nil {
		return err
	}
	defer conn.Close()

	dn := fmt.Sprintf("cn=%s,ou=people,%s", name, l.dn)
	req := ldap.NewDelRequest(dn, []ldap.Control{})
	err = conn.Del(req)

	if ldap.IsErrorWithCode(err, 32) {
		return UserNotFound
	}
	return err
}

func (l *Ldap) AddGroup(name, user string) error {
	conn, err := l.loginAdmin(l.username, l.password)
	if err != nil {
		return err
	}
	defer conn.Close()

	dn := fmt.Sprintf("cn=%s,ou=groups,%s", name, l.dn)
	req := ldap.NewAddRequest(dn)
	req.Attribute("description", []string{fmt.Sprintf("%s group", name)})
	req.Attribute("objectclass", []string{"groupofnames"})
	req.Attribute("member", []string{fmt.Sprintf("cn=%s,ou=people,%s", user, l.dn)})
	err = conn.Add(req)

	if ldap.IsErrorWithCode(err, 68) {
		return GroupAlreadyExists
	}
	return err
}

func (l *Ldap) DeleteGroup(name string) error {
	conn, err := l.loginAdmin(l.username, l.password)
	if err != nil {
		return err
	}
	defer conn.Close()

	dn := fmt.Sprintf("cn=%s,ou=groups,%s", name, l.dn)
	req := ldap.NewDelRequest(dn, []ldap.Control{})
	err = conn.Del(req)

	if ldap.IsErrorWithCode(err, 32) {
		return GroupNotFound
	}
	return err
}

func (l *Ldap) GetUsersGroups(name string) ([]string, bool) {
	conn, err := l.newConnection()
	if err != nil {
		return nil, false
	}
	defer conn.Close()

	dn := fmt.Sprintf("cn=%s,ou=people,%s", name, l.dn)

	entries, err := l.search(conn, dn, fmt.Sprintf("(&(objectclass=person)(cn=%s))", name), []string{"memberOf"})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	if len(entries) > 1 {
		return []string{}, false
	}

	users := []string{}
	entry := entries[0]
	values := entry.GetAttributeValues("memberOf")
	for _, value := range values {
		start := strings.Index(value, "cn=")
		if start == 0 {
			end := strings.Index(value, ",")
			users = append(users, value[3:end])
		}
	}

	return users, true
}

func (l *Ldap) GetUsersFromGroup(name string) ([]string, bool) {
	conn, err := l.newConnection()
	if err != nil {
		return nil, false
	}
	defer conn.Close()

	dn := fmt.Sprintf("cn=%s,ou=groups,%s", name, l.dn)

	entries, err := l.search(conn, dn, "(&(objectClass=groupofnames))", []string{"member"})
	if err != nil {
		log.Println(err)
		return nil, false
	}

	if len(entries) > 1 {
		return []string{}, false
	}

	users := []string{}
	entry := entries[0]
	values := entry.GetAttributeValues("member")
	for _, value := range values {
		start := strings.Index(value, "cn=")
		if start == 0 {
			end := strings.Index(value, ",")
			users = append(users, value[3:end])
		}
	}

	return users, true
}

func (l *Ldap) ChangeUserPassword(username, oldPassword, newPassword string) error {
	conn, err := l.loginUser(username, oldPassword)

	if err != nil {
		return err
	}

	passwordModifyRequest := ldap.NewPasswordModifyRequest("", oldPassword, newPassword)
	_, err = conn.PasswordModify(passwordModifyRequest)

	return err
}

func (l *Ldap) UpdateUser(user *models.User) error {
	if l.Validate(user.Username, user.Password) == false {
		return UserAuthFailure
	}
	conn, err := l.loginAdmin(l.username, l.password)
	if err != nil {
		return err
	}
	defer conn.Close()

	modify := ldap.NewModifyRequest(fmt.Sprintf("cn=%s,ou=people,%s", user.Username, l.dn))
	modify.Replace("mail", []string{user.Email})
	modify.Replace("sn", []string{user.LastName})
	modify.Replace("givenName", []string{user.FirstName})

	err = conn.Modify(modify)
	return err
}

func (l *Ldap) Validate(username, password string) bool {
	conn, err := l.loginUser(username, password)
	if err != nil {
		return false
	}

	conn.Close()
	return true
}

func (l *Ldap) loginAdmin(username, password string) (*ldap.Conn, error) {
	return l.login(fmt.Sprintf("cn=%s,%s", username, l.dn), password)
}

func (l *Ldap) loginUser(username, password string) (*ldap.Conn, error) {
	return l.login(fmt.Sprintf("cn=%s,ou=people,%s", username, l.dn), password)
}

func (l *Ldap) login(cn, password string) (*ldap.Conn, error) {
	conn, err := l.newConnection()
	if err != nil {
		return nil, err
	}
	err = conn.Bind(cn, password)
	if err != nil {
		conn.Close()
		return nil, UserAuthFailure
	}

	return conn, nil
}

func (l *Ldap) createTLD(name string) error {
	conn, err := l.loginAdmin(l.username, l.password)
	if err != nil {
		return err
	}
	defer conn.Close()

	req := ldap.NewAddRequest(fmt.Sprintf("ou=%s,%s", name, l.dn))
	req.Attribute("description", []string{"Group hierarchy"})
	req.Attribute("objectclass", []string{"organizationalunit"})
	err = conn.Add(req)

	if ldap.IsErrorWithCode(err, 68) {
		return GroupAlreadyExists
	}
	return err
}

func (l *Ldap) search(conn *ldap.Conn, dn, query string, attrs []string) ([]*ldap.Entry, error) {
	req := ldap.NewSearchRequest(dn,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		query,
		attrs, // A list attributes to retrieve
		nil)
	res, err := conn.Search(req)
	if err != nil {
		return nil, err
	}

	return res.Entries, nil
}
