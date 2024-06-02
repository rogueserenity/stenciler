Feature: Simple Processing

  Scenario: Processing a template that contains only raw copy files
    Given I have template with only raw copy files
    When I run stenciler init with the template in an empty directory
    Then I see the current directory initialized with the template data
