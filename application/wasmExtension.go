package application

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/CosmWasm/wasmd/x/wasm"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/persistenceOne/persistenceSDK/constants"
	"github.com/persistenceOne/persistenceSDK/modules/assets/transactions/mint"
	"github.com/persistenceOne/persistenceSDK/schema/types/base"
)

// this is for adding raw messages to wasm //

type customMessage struct {
	MsgType string          `json:"msgtype,required"`
	Raw     json.RawMessage `json:"raw,omitempty"`
}

// Type will be assets/mint , assets/burn, assets/Mmtate , (like codec register types)

func WasmCustomMessageEncoder(codec *codec.Codec) *wasm.MessageEncoders {
	return &wasm.MessageEncoders{
		Custom: customEncoder(codec),
	}
}

func customEncoder(codec *codec.Codec) wasm.CustomEncoder {
	return func(sender sdkTypes.AccAddress, rawMessage json.RawMessage) ([]sdkTypes.Msg, error) {
		var customMessage customMessage
		err := json.Unmarshal(rawMessage, &customMessage)
		if err != nil {
			return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
		}

		switch customMessage.MsgType {
		case "assets/mint":
			return assetsMintEncoder(codec, sender, customMessage.Raw)
		case "assets/mutate":
			return assetsMutateEncoder(codec, sender, customMessage.Raw)
		case "assets/burn":
			return assetsBurnEncoder(codec, sender, customMessage.Raw)
		default:
			return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant not supported in SDK")
		}
	}
}

func assetsMintEncoder(_ *codec.Codec, sender sdkTypes.AccAddress, rawMessage json.RawMessage) ([]sdkTypes.Msg, error) {
	if rawMessage != nil {
		var assetMintMessage AssetMintMessage
		err := json.Unmarshal(rawMessage, &assetMintMessage)
		if err != nil {
			return nil, sdkErrors.Wrap(sdkErrors.ErrJSONUnmarshal, err.Error())
		}
		return encodeAssetMintMessage(sender, assetMintMessage)
	}
	return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "rawMessage cannot be nil or empty")
}

func encodeAssetMintMessage(sender sdkTypes.AccAddress, assetMintMessage AssetMintMessage) ([]sdkTypes.Msg, error) {
	properties := base.ReadProperties(assetMintMessage.Properties)
	if len(properties.GetList()) > constants.MaxTraitCount {
		panic(errors.New(fmt.Sprintf("")))
	}
	mintMessage := mint.Message{
		From:             sender,
		ClassificationID: base.NewID(assetMintMessage.ClassificationID),
	}
	return []sdkTypes.Msg{mintMessage}, nil
}

func assetsMutateEncoder(_ *codec.Codec, _ sdkTypes.AccAddress, _ json.RawMessage) ([]sdkTypes.Msg, error) {
	return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant assetMutate not supported")
}

func assetsBurnEncoder(_ *codec.Codec, _ sdkTypes.AccAddress, _ json.RawMessage) ([]sdkTypes.Msg, error) {
	return nil, sdkErrors.Wrap(wasm.ErrInvalidMsg, "Custom variant assetBurn not supported")
}

// AssetMintMessage should look like rest request, or similar and should be convertible to sdk message
type AssetMintMessage struct {
	From             string `json:"from"`
	MaintainersID    string `json:"maintainersID"`
	ClassificationID string `json:"classificationID"`
	Properties       string `json:"properties"`
	Lock             int64  `json:"lock"`
	Burn             int64  `json:"burn"`
}
