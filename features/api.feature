Feature: Legacy Cache API

  Scenario: Non existing endpoint
    When I GET "/non-existing-endpoint"
    Then the HTTP status code should be "404"
