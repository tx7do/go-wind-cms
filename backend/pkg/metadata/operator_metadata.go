package metadata

import (
	"context"
	"encoding/base64"

	"github.com/go-kratos/kratos/v2/encoding"
	_ "github.com/go-kratos/kratos/v2/encoding/json"
	_ "github.com/go-kratos/kratos/v2/encoding/proto"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/transport"

	authenticationV1 "go-wind-cms/api/gen/go/authentication/service/v1"
)

const (
	mdOperatorKey  = "x-md-operator"
	mdSignatureKey = "x-md-signature"
)

var codec = encoding.GetCodec("proto")

func NewContext(ctx context.Context, info *authenticationV1.OperatorMetadata) (context.Context, error) {
	str, err := EncodeOperatorMetadata(info)
	if err != nil {
		return ctx, err
	}
	//log.Debugf("metadata: adding operator metadata to context: %s", str)
	return metadata.AppendToClientContext(ctx, mdOperatorKey, str), nil
}

func FromServerContext(ctx context.Context) (*authenticationV1.OperatorMetadata, error) {
	return FromContext(ctx, true)
}

func FromClientContext(ctx context.Context) (*authenticationV1.OperatorMetadata, error) {
	return FromContext(ctx, false)
}

func FromContext(ctx context.Context, server bool) (*authenticationV1.OperatorMetadata, error) {
	var md metadata.Metadata
	var ok bool
	if server {
		md, ok = metadata.FromServerContext(ctx)
	} else {
		md, ok = metadata.FromClientContext(ctx)
	}
	if !ok {
		return nil, ErrNoMetadata
	}

	val := md.Get(mdOperatorKey)
	if val == "" {
		log.Debugf("metadata: no operator metadata found in context [%v]", md)
		return nil, ErrNoOperatorHeader
	}

	info, err := DecodeOperatorMetadata(val)
	if err != nil {
		return nil, ErrInvalidOperator
	}

	return info, nil
}

// EncodeOperatorMetadata encodes OperatorMetadata into a base64 string
func EncodeOperatorMetadata(info *authenticationV1.OperatorMetadata) (string, error) {
	if info == nil {
		return "", ErrNilOperatorMetadata
	}
	if codec == nil {
		return "", ErrCodecNotInitialized
	}

	b, err := codec.Marshal(info)
	if err != nil {
		log.Errorf("failed to marshal operator metadata: %v", err)
		return "", err
	}
	str := base64.RawStdEncoding.EncodeToString(b)
	return str, nil
}

// DecodeOperatorMetadata decodes a base64 string into OperatorMetadata
func DecodeOperatorMetadata(str string) (*authenticationV1.OperatorMetadata, error) {
	if str == "" {
		return nil, ErrEmptyOperatorMetadataString
	}
	if codec == nil {
		return nil, ErrCodecNotInitialized
	}

	b, err := base64.RawStdEncoding.DecodeString(str)
	if err != nil {
		log.Errorf("failed to decode operator metadata: %v", err)
		return nil, err
	}

	info := &authenticationV1.OperatorMetadata{}
	if err = codec.Unmarshal(b, info); err != nil {
		log.Errorf("failed to unmarshal operator metadata: %v", err)
		return nil, err
	}
	return info, nil
}

func SetOperatorToRequestHeader(reqHeaders transport.Header, info *authenticationV1.OperatorMetadata) error {
	if reqHeaders == nil {
		return ErrNilRequestHeader
	}

	str, err := EncodeOperatorMetadata(info)
	if err != nil {
		return err
	}
	reqHeaders.Set(mdOperatorKey, str)
	return nil
}

func GetOperatorFromRequestHeader(reqHeaders transport.Header) (*authenticationV1.OperatorMetadata, error) {
	if reqHeaders == nil {
		return nil, ErrNilRequestHeader
	}

	val := reqHeaders.Get(mdOperatorKey)
	if val == "" {
		log.Debugf("metadata: no operator metadata found in request headers [%v]", reqHeaders)
		return nil, ErrNoOperatorHeader
	}

	info, err := DecodeOperatorMetadata(val)
	if err != nil {
		return nil, ErrInvalidOperator
	}

	return info, nil
}
