Feature: Simple Processing

  Scenario: Processing a template that contains only raw copy files
    Given I have a local template with only raw copy files
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  Scenario: Processing a template that a single templated file with no prompts
    Given I have a local template with a templated file and no prompts
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  @prompt
  Scenario: Processing a template that a single templated file with a prompt with no default value
    Given I have a local template with a templated file that prompts with no default value
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data
