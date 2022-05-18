Feature: Test

    Feature Description

    Scenario: Test Step functions
        Given I create iam testing role
        And I create step function state machine from "test" JSON file with "HelloWorld" name
        And I create lambda function
        And I create a DynamoDB table with "TransactionHistoryTable" name and "TransactionId" string index
        When I start step function execution with "Test" name
        Then I get step function execution describe and store response in "test-1" key
        And I validate step function output stored in "test-1" key is equal to
        """
        hello world
        """