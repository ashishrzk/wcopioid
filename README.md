#HEALTH INSURANCE CLAIM BLOCKCHAIN SOLUTION

### Use Case 1 - Out of Network Claims

* **Personas:** Insurance Payer
* **Technology:** Oracle Blockchain Cloud Service, Oracle Health Insurance *(possibly)*
* **Outline:**
This demo will show how Oracle Blockchain Cloud Service can be used to increase a payer's visiblity into claims from out of network providers. We will take the example of a batch of claims from providers with which the payer has no managed care contract. These claims are be written to provider's blockchain network.

The blockchain record ids for the claims are sent to a servicer that also has access to the blockchain, and that provider is given access to write information to these records. That servicer will then being their work on the claims, accessing them and updating them on the blockchain network.

We'll show that, since the records remain on the blockchain network, the payer will always have the ability to check the records and see their current status. Once the servicer marks their work on the claims as complete, the payer will then initiate their payments.

If possible, this demo would show a version of Oracle Health Insurance that we've customized writing and reading the data to the blockchain network. If that is not possible, then we will build a similar front end and API to read and write the records from the blockchain network.


### IOT Traceability Use Case

