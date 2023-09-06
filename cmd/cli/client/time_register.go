package client

import (
	"fmt"

	"ponto-menos/internal/punchclock"
	"ponto-menos/internal/session"
	"ponto-menos/internal/sign"
	"ponto-menos/internal/timecard"
	"strings"

	"go.uber.org/zap"

	"github.com/spf13/cobra"
)

var timeRegisterCmd = &cobra.Command{
	Use:   "time-register",
	Short: "Use this to clock in",
	Long:  `This command will clock in for you on the PontoMais platform using the current time`,
	Run: func(cmd *cobra.Command, args []string) {
		user, _ := cmd.Flags().GetString("user")
		if isEmpty(user) {
			fmt.Println("To clock in you need insert your user e-mail")
			return
		}

		password, _ := cmd.Flags().GetString("password")
		if isEmpty(password) {
			fmt.Println("To clock in you need insert your user password")
			return
		}

		credentials := sign.Credentials{User: user, Password: password}

		puncher := punchclock.PunchClock{Sign: &sign.UserSign{}, Session: &session.UserSession{}, TimeCard: &timecard.TimeCard{}}

		err := puncher.ClockIn(credentials)
		if err != nil {
			zap.L().Error(err.Error())
			fmt.Println("pontos-menos: error: an error occurred during clock in, check the system logs")
		}
	},
}

func init() {
	rootCmd.AddCommand(timeRegisterCmd)
	timeRegisterCmd.Flags().String("user", "", "Set user e-mail used in PontoMais")
	timeRegisterCmd.Flags().String("password", "", "Set password used in PontoMais")
}

func isEmpty(value string) bool {
	return value == "" || strings.TrimSpace(value) == ""
}
