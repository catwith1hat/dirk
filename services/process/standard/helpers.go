// Copyright © 2020 Attestant Limited.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package standard

import (
	context "context"

	"github.com/attestantio/dirk/core"
	"github.com/attestantio/dirk/services/checker"
	"github.com/opentracing/opentracing-go"
)

// checkAccess returns true if the client can access the account.
func (s *Service) checkAccess(ctx context.Context, credentials *checker.Credentials, accountName string, action string) core.Result {
	span, ctx := opentracing.StartSpanFromContext(ctx, "services.process.checkAccess")
	defer span.Finish()

	if s.checkerSvc.Check(ctx, credentials, accountName, action) {
		return core.ResultSucceeded
	}

	return core.ResultDenied
}
