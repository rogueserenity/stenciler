Feature: Local Template Processing

  Scenario: Processing a template with contains only raw copy files
    Given I have a local template with only raw copy files
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  Scenario: Processing a template with a single templated file with no prompts
    Given I have a local template with a templated file and no prompts
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  @prompt
  Scenario: Processing a template with a single templated file with a prompt with no default value
    Given I have a local template with a templated file that prompts with no default value
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  @prompt
  Scenario: Processing a template with a single templated file with a prompt with a default value
    Given I have a local template with a templated file that prompts with a default value
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  @prompt
  Scenario: Processing a template with a single templated file with a prompt with no default value and a hook
    Given I have a local template with a templated file that prompts with no default value and a hook
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  @prompt
  Scenario: Processing a template with a single templated file with a prompt with a default value and a hook
    Given I have a local template with a templated file that prompts with a default value and a hook
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  Scenario: Processing a template with a single pre-init hook
    Given I have a local template with a pre-init hook
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  Scenario: Processing a template with multiple pre-init hooks
    Given I have a local template with multiple pre-init hooks
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  Scenario: Processing a template with a single post-init hook
    Given I have a local template with a post-init hook
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data

  Scenario: Processing a template with multiple post-init hooks
    Given I have a local template with multiple post-init hooks
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data
