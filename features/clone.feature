Feature: Cloning a repository

  Scenario Outline: Cloning a repository from GitHub
    Given I have a "<visibility>" repository on GitHub
    And I have the "<protocol>" URL of the repository
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the repo template data

    Examples:
      | visibility | protocol |
      | public     | HTTPS    |
      | private    | HTTPS    |
      | public     | SSH      |
      | private    | SSH      |
