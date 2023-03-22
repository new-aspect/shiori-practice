package webserver

import (
	"encoding/json"
	"fmt"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/new-aspect/shiori-practice/internal/database"
	"github.com/new-aspect/shiori-practice/internal/model"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"time"
)

// apiLogin is handler for POST /api/login
func (h *handler) apiLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	// Decode request
	request := struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Remember bool   `json:"remember"`
		Owner    bool   `json:"owner"`
	}{}

	// 原来还可以这样decode
	err := json.NewDecoder(r.Body).Decode(&request)
	CheckError(err)

	// Prepare function to generate session
	genSession := func(account model.Account, expTime time.Duration) {
		// Create session ID
		sessionID, err := uuid.NewV4()
		CheckError(err)

		// Save session ID to cache
		strSessionID := sessionID.String()
		// 这里的h的SessionCache是我们自己定义的一个缓存，这个缓存是存在后端的
		h.SessionCache.Set(strSessionID, account, expTime)

		// Save user's session IDs to cache as well
		// useful for mass logout
		sessionIDs := []string{strSessionID}
		if val, found := h.UserCache.Get(request.Username); found {
			sessionIDs = val.([]string)
			sessionIDs = append(sessionIDs, strSessionID)
		}
		h.UserCache.Set(request.Username, sessionIDs, -1)

		// Send login result
		account.Password = ""
		loginResult := struct {
			Session string        `json:"session"`
			Account model.Account `json:"account"`
			Expires string        `json:"expires"`
		}{strSessionID, account, time.Now().UTC().Add(expTime).Format(time.RFC1123)}

		w.Header().Set("Content-Type", "application/json")
		err = json.NewEncoder(w).Encode(&loginResult)
		CheckError(err)
	}

	// Check if user's database is empty or there are no owner.
	// If yes, and user uses default account, let him in.
	searchOptions := database.GetAccountsOptions{
		Owner: true,
	}

	accounts, err := h.DB.GetAccounts(ctx, searchOptions)
	CheckError(err)

	if len(accounts) == 0 && request.Username == "admin" && request.Password == "admin" {
		genSession(model.Account{
			Username: "admin",
			Owner:    true,
		}, time.Hour)
		return
	}

	// Get account data form database
	account, exit, err := h.DB.GetAccount(ctx, request.Username)
	CheckError(err)

	if !exit {
		panic(fmt.Errorf("username doesn't exist"))
	}

	// Compare password with database
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(request.Password))
	if err != nil {
		panic(fmt.Errorf("username and password don't match"))
	}

	// If login request is as owner, make sure this account is owner
	if request.Owner && !account.Owner {
		panic(fmt.Sprintf("account level is not sufficient as owner"))
	}

	// Calculate expiration time
	expTime := time.Hour
	if request.Remember {
		expTime = time.Hour * 24 * 30
	}

	// Create session
	genSession(account, expTime)
}
