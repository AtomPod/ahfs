package models

import "fmt"

type ErrUserNotExist struct {
	ID       uint
	Email    string
	Username string
}

func IsErrUserNotExist(err error) bool {
	_, ok := err.(ErrUserNotExist)
	return ok
}

func (err ErrUserNotExist) Error() string {
	return fmt.Sprintf("user does not exist [id: %d, email: %s, username: %s]", err.ID, err.Email, err.Username)
}

type ErrEmailAlreadyUsed struct {
	Email string
}

func IsErrEmailAlreadyUsed(err error) bool {
	_, ok := err.(ErrEmailAlreadyUsed)
	return ok
}

func (err ErrEmailAlreadyUsed) Error() string {
	return fmt.Sprintf("email already used [email: %s]", err.Email)
}

type ErrUsernameAlreadyUsed struct {
	Username string
}

func IsErrUsernameAlreadyUsed(err error) bool {
	_, ok := err.(ErrUsernameAlreadyUsed)
	return ok
}

func (err ErrUsernameAlreadyUsed) Error() string {
	return fmt.Sprintf("username already used [username: %s]", err.Username)
}

type ErrModifyRootFile struct {
	ID    uint
	Owner uint
}

func IsErrModifyRootFile(err error) bool {
	_, ok := err.(ErrModifyRootFile)
	return ok
}

func (err ErrModifyRootFile) Error() string {
	return fmt.Sprintf("cannot modify root file [id: %d, owner: %d]", err.ID, err.Owner)
}

type ErrFileNotExist struct {
	ID     uint
	Path   string
	Owner  uint
	FileID string
}

func IsErrFileNotExist(err error) bool {
	_, ok := err.(ErrFileNotExist)
	return ok
}

func (err ErrFileNotExist) Error() string {
	return fmt.Sprintf("file does not exist [id: %d, path: %s, owner: %d, file_id: %s]", err.ID, err.Path, err.Owner, err.FileID)
}

type ErrFileAlreadyExist struct {
	ID     uint
	Path   string
	Owner  uint
	FileID string
}

func IsErrFileAlreadyExist(err error) bool {
	_, ok := err.(ErrFileAlreadyExist)
	return ok
}

func (err ErrFileAlreadyExist) Error() string {
	return fmt.Sprintf("file already exist [id: %d, path: %s, owner: %d, file_id: %s]", err.ID, err.Path, err.Owner, err.FileID)
}

type ErrFileNotDirectory struct {
	ID   uint
	Path string
}

func IsErrFileNotDirectory(err error) bool {
	_, ok := err.(ErrFileNotDirectory)
	return ok
}

func (err ErrFileNotDirectory) Error() string {
	return fmt.Sprintf("path is not a directory [id: %d, path: %s]", err.ID, err.Path)
}

type ErrFileParentNotDirectory struct {
	ID   uint
	Path string
}

func IsErrFileParentNotDirectory(err error) bool {
	_, ok := err.(ErrFileParentNotDirectory)
	return ok
}

func (err ErrFileParentNotDirectory) Error() string {
	return fmt.Sprintf("parent is not a directory [id: %d, path: %s]", err.ID, err.Path)
}

type ErrFileLocked struct {
	ID uint
}

func IsErrFileLocked(err error) bool {
	_, ok := err.(ErrFileLocked)
	return ok
}

func (err ErrFileLocked) Error() string {
	return fmt.Sprintf("file is locked [id: %d]", err.ID)
}

type ErrFileUnlockFailed struct {
	ID     uint
	LockID string
}

func IsFileUnlockFailed(err error) bool {
	_, ok := err.(ErrFileUnlockFailed)
	return ok
}

func (err ErrFileUnlockFailed) Error() string {
	return fmt.Sprintf("file cannot unlock, maybe timeout or lock id is incorret [id: %d, lockID: %s]", err.ID, err.LockID)
}

type ErrUserMaxFileCapacityLimit struct {
	UserID uint
}

func IsErrFileMaxSizeLimit(err error) bool {
	_, ok := err.(ErrUserMaxFileCapacityLimit)
	return ok
}

func (err ErrUserMaxFileCapacityLimit) Error() string {
	return fmt.Sprintf("user file capacity is fulled [id: %d]", err.UserID)
}

type ErrOAuth2TokenNotExist struct {
	Code         string
	AccessToken  string
	RefreshToken string
}

func IsErrOAuth2TokenNotExist(err error) bool {
	_, ok := err.(ErrOAuth2TokenNotExist)
	return ok
}
func (e ErrOAuth2TokenNotExist) Error() string {
	return fmt.Sprintf("oauth2 token does not exist [code: %s, access_token: %s, refresh_token: %s]", e.Code, e.AccessToken, e.RefreshToken)
}

type ErrOAuth2ApplicationNotExist struct {
	ID       uint
	ClientID string
	Name     string
	UserID   uint
}

func IsErrOAuth2ApplicationNotExist(err error) bool {
	_, ok := err.(ErrOAuth2ApplicationNotExist)
	return ok
}

func (e ErrOAuth2ApplicationNotExist) Error() string {
	return fmt.Sprintf("oauth2 application does not exist [id: %d, clientID: %s, name: %s, user_id: %d]", e.ID, e.ClientID, e.Name, e.UserID)
}

type ErrAuthTokenNotExist struct {
	Code   string
	ID     uint
	UserID uint
}

func IsErrAuthTokenNotExist(err error) bool {
	_, ok := err.(ErrAuthTokenNotExist)
	return ok
}

func (e ErrAuthTokenNotExist) Error() string {
	return fmt.Sprintf("auth token does not exist [id: %d, user_id: %d, code: %s]", e.ID, e.UserID, e.Code)
}

type ErrAuthTokenEmpty struct {
}

func IsErrAuthTokenEmpty(err error) bool {
	_, ok := err.(ErrAuthTokenEmpty)
	return ok
}

func (e ErrAuthTokenEmpty) Error() string {
	return fmt.Sprintf("auth token cannot be empty")
}
