Feature: Confirm

    Scenario: With help and receive a yes
        Given I see a confirm prompt "Confirm? [? for help] (y/N)", I ask for help and see "This action cannot be undone"
        And then I see another confirm prompt "Confirm? (y/N)", I answer yes

        Then ask for confirm "Confirm?" with help "This action cannot be undone", receive yes

    Scenario: With help and receive a no
        Given I see a confirm prompt "Confirm? [? for help] (y/N)", I ask for help and see "This action cannot be undone"
        And then I see another confirm prompt "Confirm? (y/N)", I answer no

        Then ask for confirm "Confirm?" with help "This action cannot be undone", receive no

    Scenario: Without help and receive a yes
        Given I see a confirm prompt "Confirm? (y/N)", I answer yes

        Then ask for confirm "Confirm?", receive yes

    Scenario: Without help and receive a no
        Given I see a confirm prompt "Confirm? (y/N)", I answer no

        Then ask for confirm "Confirm?", receive no

    Scenario: Interrupted
        Given I see a confirm prompt "Confirm? (y/N)", I interrupt

        Then ask for confirm "Confirm?", get interrupted

    Scenario: Invalid answer
        Given I see a confirm prompt "Confirm? (y/N)", I answer "nahhh"
        And then I see another confirm prompt "Confirm? (y/N)", I answer no

        Then ask for confirm "Confirm?", receive no
