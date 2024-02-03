Feature: Upsert Cache Time

  Scenario: Create Cache Time resource
    Given the document with "_id" set to "5d41402abc4b2a76b9719d911017c592" does not exist in the "cachetimes" collection
    When I PUT "/v1/cache-times/5d41402abc4b2a76b9719d911017c592"
            """
            {
                "path": "/my-path",
                "etag": "test-etag",
                "collection_id": 123456,
                "release_time": "2024-01-31T01:23:45.678Z"
            }
            """
    Then the HTTP status code should be "204"

  Scenario: Update Cache Time resource
    Given the following document exists in the "cachetimes" collection:
            """
            {
                "_id": "5d41402abc4b2a76b9719d911017c592",
                "path": "/my-path",
                "etag": "test-etag",
                "collection_id": 123456,
                "release_time": "2024-01-31T01:23:45.678Z"
            }
            """
    When I PUT "/v1/cache-times/5d41402abc4b2a76b9719d911017c592"
            """
            {
                "path": "/some/other/path",
                "etag": "a-different-etag",
                "collection_id": 999,
                "release_time": "1999-12-23T11:22:33.444Z"
            }
            """
    Then the HTTP status code should be "204"

  Scenario: Upsert Cache Time resource with empty body
    Given the document with "_id" set to "5d41402abc4b2a76b9719d911017c592" does not exist in the "cachetimes" collection
    When I PUT "/v1/cache-times/5d41402abc4b2a76b9719d911017c592"
            """
            """
    Then the HTTP status code should be "400"