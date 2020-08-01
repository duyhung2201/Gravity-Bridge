package cli

import (
	"bufio"
	"crypto/ecdsa"
	"fmt"
	"log"

	ethCrypto "github.com/ethereum/go-ethereum/crypto"
	"github.com/spf13/cobra"

	"github.com/althea-net/peggy/module/x/nameservice/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	nameserviceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Nameservice transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	nameserviceTxCmd.AddCommand(flags.PostCommands(
		// GetCmdBuyName(cdc),
		// GetCmdSetName(cdc),
		// GetCmdDeleteName(cdc),
		CmdUpdateEthAddress(cdc),
		CmdValsetRequest(cdc),
		CmdValsetConfirm(storeKey, cdc),
	)...)

	return nameserviceTxCmd
}

// GetCmdUpdateEthAddress updates the network about the eth address that you have on record.
func CmdUpdateEthAddress(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "update-eth-addr [eth private key]",
		Short: "update your eth address which will be used for peggy if you are a validator",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			cosmosAddr := cliCtx.GetFromAddress()

			privKeyString := args[0][2:]

			// Make Eth Signature over validator address
			privateKey, err := ethCrypto.HexToECDSA(privKeyString)
			if err != nil {
				log.Fatal(err)
			}

			hash := ethCrypto.Keccak256Hash(cosmosAddr) // TODO: Can probably skip the "Hash" struct and use ethCrypto.Keccak256
			signature, err := ethCrypto.Sign(hash.Bytes(), privateKey)
			if err != nil {
				log.Fatal(err)
			}

			// You've got to do all this to get an Eth address from the private key
			publicKey := privateKey.Public()
			publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
			if !ok {
				log.Fatal("error casting public key to ECDSA")
			}
			ethAddress := ethCrypto.PubkeyToAddress(*publicKeyECDSA).Hex()

			// Make the message
			msg := types.NewMsgSetEthAddress(ethAddress, cosmosAddr, signature)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Send it
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func CmdValsetRequest(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-request",
		Short: "request that the validators sign over the current valset",
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))
			cosmosAddr := cliCtx.GetFromAddress()

			// Make the message
			msg := types.NewMsgValsetRequest(cosmosAddr)

			// Send it
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func CmdValsetConfirm(storeKey string, cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "valset-confirm [nonce] [eth private key]",
		Short: "this is used by validators to sign a valset with a particular nonce if it exists",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			inBuf := bufio.NewReader(cmd.InOrStdin())
			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

			nonce := args[0]

			res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/valsetRequest/%s", storeKey, nonce), nil)
			if err != nil {
				fmt.Printf("could not get valset")
				return nil
			}

			var valset types.Valset
			cdc.MustUnmarshalJSON(res, &valset)
			checkpoint := valset.GetCheckpoint()

			// Make Eth Signature over valset
			privKeyString := args[0][2:]
			privateKey, err := ethCrypto.HexToECDSA(privKeyString)
			if err != nil {
				log.Fatal(err)
			}
			signature, err := ethCrypto.Sign(checkpoint, privateKey)
			if err != nil {
				log.Fatal(err)
			}

			cosmosAddr := cliCtx.GetFromAddress()

			// Make the message
			msg := types.NewMsgValsetConfirm(valset.Nonce, cosmosAddr, signature)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			// Send it
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

// // GetCmdBuyName is the CLI command for sending a BuyName transaction
// func GetCmdBuyName(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "buy-name [name] [amount]",
// 		Short: "bid for existing name or claim new name",
// 		Args:  cobra.ExactArgs(2),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)

// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

// 			coins, err := sdk.ParseCoins(args[1])
// 			if err != nil {
// 				return err
// 			}

// 			msg := types.NewMsgBuyName(args[0], coins, cliCtx.GetFromAddress())
// 			err = msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }

// // GetCmdSetName is the CLI command for sending a SetName transaction
// func GetCmdSetName(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "set-name [name] [value]",
// 		Short: "set the value associated with a name that you own",
// 		Args:  cobra.ExactArgs(2),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

// 			// if err := cliCtx.EnsureAccountExists(); err != nil {
// 			// 	return err
// 			// }

// 			msg := types.NewMsgSetName(args[0], args[1], cliCtx.GetFromAddress())
// 			err := msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

// 			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }

// // GetCmdDeleteName is the CLI command for sending a DeleteName transaction
// func GetCmdDeleteName(cdc *codec.Codec) *cobra.Command {
// 	return &cobra.Command{
// 		Use:   "delete-name [name]",
// 		Short: "delete the name that you own along with it's associated fields",
// 		Args:  cobra.ExactArgs(1),
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			cliCtx := context.NewCLIContext().WithCodec(cdc)
// 			inBuf := bufio.NewReader(cmd.InOrStdin())
// 			txBldr := auth.NewTxBuilderFromCLI(inBuf).WithTxEncoder(utils.GetTxEncoder(cdc))

// 			msg := types.NewMsgDeleteName(args[0], cliCtx.GetFromAddress())
// 			err := msg.ValidateBasic()
// 			if err != nil {
// 				return err
// 			}

// 			// return utils.CompleteAndBroadcastTxCLI(txBldr, cliCtx, msgs)
// 			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
// 		},
// 	}
// }
