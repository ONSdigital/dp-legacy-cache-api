@GetCacheTimes
Feature: Get paginated cache times

  Scenario: Read existing Cache Time resources
    Given the following document exists in the "cachetimes" collection:
      """
      {
        "_id": "5d41402abc4b2a76b9719d911017c592",
        "path": "/my-path",
        "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
        "release_time": "2024-01-31T01:23:45.678Z"
      }
      """
    When I GET "/v1/cache-times"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 1,
        "items": [
          {
            "_id": "5d41402abc4b2a76b9719d911017c592",
            "path": "/my-path",
            "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
            "release_time": "2024-01-31T01:23:45.678Z"
          }
        ],
        "limit": 10,
        "offset": 0,
        "total_count": 1
      }
      """

  Scenario: Read existing Cache Time resources - pagination
    Given the following document exists in the "cachetimes" collection:
      """
      {
        "_id": "5d41402abc4b2a76b9719d911017c592",
        "path": "/my-path",
        "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
        "release_time": "2024-01-31T01:23:45.678Z"
      }
      """
    And the following document exists in the "cachetimes" collection:
      """
      {
        "_id": "7e57d0042b97b6f99b5e6a8d6e0ae5ae",
        "path": "/my-path2",
        "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
        "release_time": "2024-02-01T12:00:00.000Z"
      }
      """
    When I GET "/v1/cache-times"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 2,
        "items": [
          {
            "_id": "5d41402abc4b2a76b9719d911017c592",
            "path": "/my-path",
            "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
            "release_time": "2024-01-31T01:23:45.678Z"
          },
          {
            "_id": "7e57d0042b97b6f99b5e6a8d6e0ae5ae",
            "path": "/my-path2",
            "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
            "release_time": "2024-02-01T12:00:00Z"
          }
        ],
        "limit": 10,
        "offset": 0,
        "total_count": 2
      }
      """
    When I GET "/v1/cache-times?offset=1"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 1,
        "items": [
          {
            "_id": "7e57d0042b97b6f99b5e6a8d6e0ae5ae",
            "path": "/my-path2",
            "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
            "release_time": "2024-02-01T12:00:00Z"
          }
        ],
        "limit": 10,
        "offset": 1,
        "total_count": 2
      }
      """
    When I GET "/v1/cache-times?limit=1"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 1,
        "items": [
          {
            "_id": "5d41402abc4b2a76b9719d911017c592",
            "path": "/my-path",
            "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
            "release_time": "2024-01-31T01:23:45.678Z"
          }
        ],
        "limit": 1,
        "offset": 0,
        "total_count": 2
      }
      """

  Scenario: Read existing Cache Time resources - filter by release time
    Given I am authorised
    # Having to PUT here instead of inserting directly as direct insertion
    # does not store the release_time as a string
    And I PUT "/v1/cache-times/5d41402abc4b2a76b9719d911017c592"
      """
      {
        "path": "/my-path",
        "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
        "release_time": "2024-01-31T01:23:45.678Z"
      }
      """
    And I PUT "/v1/cache-times/7e57d0042b97b6f99b5e6a8d6e0ae5ae"
      """
      {
        "path": "/my-path2",
        "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
        "release_time": "2024-02-01T12:00:00.000Z"
      }
      """
    When I GET "/v1/cache-times?release_time=2024-01-31T01:23:45.678Z"
    Then I should receive the following JSON response with status "200":
      """
      {
        "count": 1,
        "items": [
          {
            "_id": "5d41402abc4b2a76b9719d911017c592",
            "path": "/my-path",
            "collection_id": "test-1a19e3462937d85804752375daa00ba41d1b6625d396f21000e3c4571ebf2606",
            "release_time": "2024-01-31T01:23:45.678Z"
          }
        ],
        "limit": 10,
        "offset": 0,
        "total_count": 1
      }
      """
