package game

import (
	"context"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/matryer/is"
	"github.com/rs/zerolog/log"

	"github.com/domino14/macondo/board"
	macondoconfig "github.com/domino14/macondo/config"
	macondogame "github.com/domino14/macondo/game"
	macondopb "github.com/domino14/macondo/gen/api/proto/macondo"

	"github.com/domino14/liwords/pkg/config"
	"github.com/domino14/liwords/pkg/entity"
	"github.com/domino14/liwords/pkg/stores/common"
	"github.com/domino14/liwords/pkg/stores/user"
	pkguser "github.com/domino14/liwords/pkg/user"
	pb "github.com/domino14/liwords/rpc/api/proto/ipc"
)

var DefaultConfig = macondoconfig.Config{
	StrategyParamsPath:        os.Getenv("STRATEGY_PARAMS_PATH"),
	LexiconPath:               os.Getenv("LEXICON_PATH"),
	LetterDistributionPath:    os.Getenv("LETTER_DISTRIBUTION_PATH"),
	DefaultLexicon:            "NWL18",
	DefaultLetterDistribution: "English",
}

func newMacondoGame(users [2]*entity.User) *macondogame.Game {
	rules, err := macondogame.NewBasicGameRules(
		&DefaultConfig, DefaultConfig.DefaultLexicon,
		board.CrosswordGameLayout, DefaultConfig.DefaultLetterDistribution,
		macondogame.CrossScoreOnly, "")
	if err != nil {
		panic(err)
	}
	var players []*macondopb.PlayerInfo

	for _, u := range users {
		players = append(players, &macondopb.PlayerInfo{
			Nickname: u.Username,
			UserId:   u.UUID,
			RealName: u.RealName(),
		})
	}

	mcg, err := macondogame.NewGame(rules, players)
	if err != nil {
		panic(err)
	}
	return mcg
}

func userStore(dbURL string) pkguser.Store {
	pool, err := common.OpenTestingDB()
	if err != nil {
		panic(err)
	}
	ustore, err := user.NewDBStore(pool)
	if err != nil {
		log.Fatal().Err(err).Msg("error")
	}
	return ustore
}

func recreateDB() pkguser.Store {
	err := common.RecreateTestDB()
	if err != nil {
		panic(err)
	}

	// Crete a user table. Initialize the user store.
	ustore := userStore(common.TestingPostgresConnDSN())
	// Insert a couple of users into the table.

	for _, u := range []*entity.User{
		{Username: "cesar", Email: "cesar@woogles.io", UUID: "mozEwaVMvTfUA2oxZfYN8k"},
		{Username: "mina", Email: "mina@gmail.com", UUID: "iW7AaqNJDuaxgcYnrFfcJF"},
		{Username: "jesse", Email: "jesse@woogles.io", UUID: "3xpEkpRAy3AizbVmDg3kdi"},
	} {
		err = ustore.New(context.Background(), u)
		if err != nil {
			log.Fatal().Err(err).Msg("error")
		}
	}

	addfakeGames(ustore)
	return ustore
}

func addfakeGames(ustore pkguser.Store) {
	protocts, err := ioutil.ReadFile("./testdata/game1/history.pb")
	if err != nil {
		log.Fatal().Err(err).Msg("error")
	}
	protocts, err = hex.DecodeString(string(protocts))
	if err != nil {
		log.Fatal().Err(err).Msg("error")
	}

	req, err := hex.DecodeString("12180a0d43726f7373776f726447616d651207656e676c697368183c2803")
	if err != nil {
		log.Fatal().Err(err).Msg("error")
	}
	// Add some fake games to the table
	store, err := NewDBStore(&config.Config{
		DBConnDSN: common.TestingPostgresConnDSN()}, ustore)
	if err != nil {
		log.Fatal().Err(err).Msg("error")
	}
	db := store.db.Exec("INSERT INTO games(created_at, updated_at, uuid, "+
		"player0_id, player1_id, timers, started, game_end_reason, winner_idx, loser_idx, "+
		"request, history, quickdata) VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)",
		"2020-07-27 04:33:45.938304+00", "2020-07-27 04:33:45.938304+00",
		"wJxURccCgSAPivUvj4QdYL", 2, 1,
		`{"lu": 1595824425928, "mo": 0, "tr": [60000, 60000], "ts": 1595824425928}`,
		true, 0, 0, 0, req, protocts,
		`{"pi":[{"nickname":"mina","rating":"1600?"},{"nickname":"cesar","rating":"500?"}]}`)

	if db.Error != nil {
		log.Fatal().Err(db.Error).Msg("error")
	}
	store.Disconnect()
	ustore.(*user.DBStore).Disconnect()

}

func teardown() {
	err := common.TeardownTestDB()
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {

	code := m.Run()
	//teardown()
	os.Exit(code)
}

func createGame(p0, p1 string, initTime int32, is *is.I) *entity.Game {
	ustore := userStore(common.TestingPostgresConnDSN())
	store, err := NewDBStore(&config.Config{
		DBConnDSN: common.TestingPostgresConnDSN()}, ustore)
	is.NoErr(err)

	u1, err := ustore.Get(context.Background(), p0)
	is.NoErr(err)

	u2, err := ustore.Get(context.Background(), p1)
	is.NoErr(err)

	mcg := newMacondoGame([2]*entity.User{u1, u2})
	mcg.StartGame()
	mcg.SetChallengeRule(macondopb.ChallengeRule_FIVE_POINT)
	mcg.SetBackupMode(macondogame.InteractiveGameplayMode)
	mcg.SetStateStackLength(1)
	entGame := entity.NewGame(mcg, &pb.GameRequest{
		InitialTimeSeconds: initTime,
		ChallengeRule:      macondopb.ChallengeRule_FIVE_POINT,
		Rules: &pb.GameRules{
			BoardLayoutName:        "CrosswordGame",
			LetterDistributionName: "english",
		},
	})
	entGame.PlayerDBIDs = [2]uint{u1.ID, u2.ID}
	entGame.Quickdata = &entity.Quickdata{
		PlayerInfo: []*pb.PlayerInfo{
			{Nickname: u1.Username, Rating: "1500?"},
			{Nickname: u2.Username, Rating: "1500?"},
		},
	}
	entGame.ResetTimersAndStart()

	err = store.Create(context.Background(), entGame)
	is.NoErr(err)

	// Clean up connections
	ustore.(*user.DBStore).Disconnect()
	store.Disconnect()

	return entGame
}

func TestCreate(t *testing.T) {
	log.Info().Msg("TestCreate")
	recreateDB()
	is := is.New(t)
	entGame := createGame("cesar", "mina", int32(60), is)

	is.True(entGame.Quickdata != nil)

	ustore := userStore(common.TestingPostgresConnDSN())
	store, err := NewDBStore(&config.Config{
		MacondoConfig: DefaultConfig,
		DBConnDSN:     common.TestingPostgresConnDSN(),
	}, ustore)
	is.NoErr(err)
	// Make sure we can fetch the game from the DB.
	log.Debug().Str("entGameID", entGame.GameID()).Msg("trying-to-fetch")
	cpy, err := store.Get(context.Background(), entGame.GameID())
	is.NoErr(err)
	is.True(cpy.Quickdata != nil)

	// Clean up connections
	ustore.(*user.DBStore).Disconnect()
	store.Disconnect()

}

func TestSet(t *testing.T) {
	log.Info().Msg("TestSet")
	recreateDB()

	is := is.New(t)
	ustore := userStore(common.TestingPostgresConnDSN())
	store, err := NewDBStore(&config.Config{
		MacondoConfig: DefaultConfig,
		DBConnDSN:     common.TestingPostgresConnDSN()}, ustore)
	is.NoErr(err)

	// Fetch the game from the backend.
	entGame, err := store.Get(context.Background(), "wJxURccCgSAPivUvj4QdYL")
	is.NoErr(err)
	// Make some modification.

	log.Debug().Interface("history", entGame.History()).Msg("test-set")
	is.Equal(entGame.History().ChallengeRule, macondopb.ChallengeRule_FIVE_POINT)
	is.Equal(entGame.NickOnTurn(), "mina")
	is.Equal(entGame.Turn(), 0)

	_, err = entGame.PlayScoringMove("8E", "AGUE", true)
	is.NoErr(err)
	// Save it back
	err = store.Set(context.Background(), entGame)
	is.NoErr(err)

	// Now, fetch the game again and see if things have updated.
	g, err := store.Get(context.Background(), "wJxURccCgSAPivUvj4QdYL")
	is.NoErr(err)
	is.Equal(g.Turn(), 1)
	is.Equal(g.NickOnTurn(), "cesar")
	// cesar is player 0, mina is player 1, but mina went first because of
	// the second_went_first flag in the history.
	is.Equal(g.RackLettersFor(0), "AEIJVVW")
	// AGUE was worth 10. The spread for cesar is therefore -10.
	is.Equal(g.SpreadFor(0), -10)

	// Clean up connections
	ustore.(*user.DBStore).Disconnect()
	store.Disconnect()
}

func TestGet(t *testing.T) {
	log.Info().Msg("TestGet")
	recreateDB()

	is := is.New(t)

	ustore := userStore(common.TestingPostgresConnDSN())
	store, err := NewDBStore(&config.Config{
		MacondoConfig: DefaultConfig,
		DBConnDSN:     common.TestingPostgresConnDSN(),
	}, ustore)
	is.NoErr(err)

	entGame, err := store.Get(context.Background(), "wJxURccCgSAPivUvj4QdYL")
	is.NoErr(err)
	log.Info().Interface("entGame history", entGame.History()).Msg("history")

	mina, err := ustore.Get(context.Background(), "mina")
	is.NoErr(err)
	log.Debug().Interface("mina", mina).Msg("playerinfo")

	is.Equal(entGame.GameID(), "wJxURccCgSAPivUvj4QdYL")
	is.Equal(entGame.NickOnTurn(), "mina")
	is.Equal(entGame.PlayerIDOnTurn(), mina.UUID)
	is.Equal(entGame.RackLettersFor(0), "AEIJVVW")
	is.Equal(entGame.RackLettersFor(1), "AEEGOUU")
	is.Equal(entGame.ChallengeRule(), macondopb.ChallengeRule_FIVE_POINT)
	is.Equal(entGame.History().ChallengeRule, macondopb.ChallengeRule_FIVE_POINT)
	// Clean up connections
	ustore.(*user.DBStore).Disconnect()
	store.Disconnect()
}

func TestListActive(t *testing.T) {
	log.Info().Msg("TestListActive")
	recreateDB()
	is := is.New(t)
	createGame("cesar", "jesse", int32(120), is)
	createGame("jesse", "mina", int32(240), is)
	ustore := userStore(common.TestingPostgresConnDSN())

	// There should be an additional game, so 3 total, from recreateDB()
	// The first game is cesar vs mina. (see TestGet)
	store, err := NewDBStore(&config.Config{
		MacondoConfig: DefaultConfig,
		DBConnDSN:     common.TestingPostgresConnDSN(),
	}, ustore)
	is.NoErr(err)

	games, err := store.ListActive(context.Background(), "")
	is.NoErr(err)
	is.Equal(len(games.GameInfo), 3)
	is.Equal(games.GameInfo[0].Players, []*pb.PlayerInfo{
		{Rating: "1600?", Nickname: "mina"},
		{Rating: "500?", Nickname: "cesar"},
	})
	is.Equal(games.GameInfo[1].Players, []*pb.PlayerInfo{
		{Rating: "1500?", Nickname: "cesar"},
		{Rating: "1500?", Nickname: "jesse"},
	})
	is.Equal(games.GameInfo[2].Players, []*pb.PlayerInfo{
		{Rating: "1500?", Nickname: "jesse"},
		{Rating: "1500?", Nickname: "mina"},
	})

	is.Equal(games.GameInfo[1].GameRequest.InitialTimeSeconds, int32(120))
	is.Equal(games.GameInfo[1].GameRequest.ChallengeRule, macondopb.ChallengeRule_FIVE_POINT)
	is.Equal(games.GameInfo[1].GameRequest.Rules, &pb.GameRules{
		BoardLayoutName:        "CrosswordGame",
		LetterDistributionName: "english",
	})
	ustore.(*user.DBStore).Disconnect()
	store.Disconnect()
}
