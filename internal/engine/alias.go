// Copyright (C) 2023 Tencent Cloud.
// Permission is hereby granted, free of charge, to any person obtaining a copy of
// this software and associated documentation files (the vectordb-sdk-java), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is furnished
// to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package engine

import (
	"context"

	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity"
	"git.woa.com/cloud_nosql/vectordb/vectordatabase-sdk-go/entity/api/alias"
)

var _ entity.AliasInterface = &implementerAlias{}

type implementerAlias struct {
	entity.SdkClient
	database entity.Database
}

func (i *implementerAlias) SetAlias(ctx context.Context, collectionName, aliasName string, option ...*entity.SetAliasOption) (*entity.SetAliasResult, error) {
	if i.database.IsAIDatabase() {
		return nil, entity.AIDbTypeError
	}
	req := new(alias.SetReq)
	req.Database = i.database.DatabaseName
	req.Collection = collectionName
	req.Alias = aliasName
	res := new(alias.SetRes)

	result := new(entity.SetAliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}

func (i *implementerAlias) DeleteAlias(ctx context.Context, aliasName string, option ...*entity.DeleteAliasOption) (*entity.DeleteAliasResult, error) {
	if i.database.IsAIDatabase() {
		return nil, entity.AIDbTypeError
	}
	req := new(alias.DeleteReq)
	req.Database = i.database.DatabaseName
	req.Alias = aliasName
	res := new(alias.DeleteRes)

	result := new(entity.DeleteAliasResult)
	err := i.Request(ctx, req, &res)
	if err != nil {
		return result, err
	}
	result.AffectedCount = res.AffectedCount
	return result, nil
}
