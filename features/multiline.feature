Feature: Multiline

    Scenario: Receive an empty answer
        Given I see a multiline prompt "Enter comment", I answer ""

        Then ask for multiline "Enter comment", receive:
        """
        """

    Scenario: Receive a multiline answer
        Given I see a multiline prompt "Enter comment", I answer:
        """
        This is the first
        line

        this is the second line
        """

        Then ask for multiline "Enter comment", receive:
        """
        This is the first
        line

        this is the second line
        """

    Scenario: Interrupted
        Given I see a multiline prompt "Enter comment", I interrupt

        Then ask for multiline "Enter comment", get interrupted
