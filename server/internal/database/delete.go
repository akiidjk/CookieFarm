package database

import "github.com/ByteTheCookies/cookieserver/internal/logger"

const (
	queryDeleteFlag = `DELETE FROM flags WHERE flag_code = ?`
)

func DeleteFlag(flag string) error {
	_, err := DB.Exec(queryDeleteFlag, flag)
	if err != nil {
		logger.Log.Error().Err(err).Msg("Failed to delete flag from DB")
		return err
	}
	logger.Log.Info().Str("flag", flag).Msg("Flag deleted successfully")
	return nil
}
