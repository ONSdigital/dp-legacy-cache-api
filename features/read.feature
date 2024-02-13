Feature: Read Cache Time

  Scenario: Read existing Cache Time resource
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
    When I GET "/v1/cache-times/5d41402abc4b2a76b9719d911017c592"
    Then I should receive the following JSON response with status "200":
      """
      {
        "_id": "5d41402abc4b2a76b9719d911017c592",
        "path": "/my-path",
        "etag": "test-etag",
        "collection_id": 123456,
        "release_time": "2024-01-31T01:23:45.678Z"
      }
      """

  Scenario: Read non-existing Cache Time resource
    When I GET "/v1/cache-times/5d41402abc4b2a76b9719d911017c592"
    Then the HTTP status code should be "404"

  Scenario: Read Cache Time resource with invalid ID format
    When I GET "/v1/cache-times/INVALID-ID"
    Then I should receive the following JSON response with status "400":
      """
      {
        "error": "validation errors: [id should be 32 characters in length, id is not lowercase, id is not a valid hexadecimal]"
      }
      """
