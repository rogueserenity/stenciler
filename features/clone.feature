@remote
Feature: Cloning a repository

  Scenario: Cloning a public repository from GitHub using HTTPS
    Given I have a "public" repository on GitHub
    And I have the "HTTPS" URL of the repository
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the repo template data

  Scenario: Cloning a public repository from GitHub using SSH
    Given I have a "public" repository on GitHub
    And I have the "SSH" URL of the repository
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the repo template data

  Scenario: Cloning a private repository from GitHub using HTTPS
    Given I have a "private" repository on GitHub
    And I have the "HTTPS" URL of the repository
    And I have a valid GitHub Authentication Token
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the repo template data

  Scenario: Cloning a private repository from GitHub using SSH
    Given I have a "private" repository on GitHub
    And I have the "SSH" URL of the repository
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the repo template data
