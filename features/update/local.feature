@update
Feature: Local Template Processing

  Scenario: Processing a template with contains only raw copy files
    Given I have a local updated template with only raw copy files
    When I run stenciler update in the current directory
    Then I see the current directory updated with the template data

  Scenario: Processing a template with a single templated file with no prompts
    Given I have a local updated template with a templated file and no prompts
    When I run stenciler update in the current directory
    Then I see the current directory updated with the template data

  Scenario: Processing a template with a single pre-update hook
    Given I have a local template with a pre-update hook
    When I run stenciler update in the current directory
    Then I see the current directory updated with the template data

  Scenario: Processing a template with multiple pre-update hooks
    Given I have a local template with multiple pre-update hooks
    When I run stenciler update in the current directory
    Then I see the current directory updated with the template data

  Scenario: Processing a template with a single post-update hook
    Given I have a local template with a post-update hook
    When I run stenciler update in the current directory
    Then I see the current directory updated with the template data

  Scenario: Processing a template with multiple post-update hooks
    Given I have a local template with multiple post-update hooks
    When I run stenciler update in the current directory
    Then I see the current directory updated with the template data

  @prompt
  Scenario: Processing a template update with a prompt but already has a value
    Given I have a local updated template with existing values
    When I run stenciler update in the current directory
    Then I see the current directory updated with the template data

  @prompt
  Scenario: Processing a template update with a prompt
    Given I have a local updated template with a new param
    When I run stenciler update in the current directory
    Then I see the current directory updated with the template data

  Scenario: Init only files are not copied during an update
    Given I have a local updated template with init-only files
    When I run stenciler update in the current directory
    Then I see the current directory updated with the template data
