syntax = "proto3";

package demo.v1;

import "google/api/annotations.proto";

// blog service is a blog demo
service BlogService {

	rpc GetArticles(GetArticlesReq) returns (GetArticlesResp) {
		option (google.api.http) = {
			get: "/v1/articles"
			additional_bindings {
				get: "/v1/author/{author_id}/articles"
			}
		};
	}
}