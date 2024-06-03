Feature: Simple Processing

  Scenario: Processing a template that contains only raw copy files
    Given I have a local template with only raw copy files
    When I run stenciler init with the repository URL in an empty directory
    Then I see the current directory initialized with the template data
