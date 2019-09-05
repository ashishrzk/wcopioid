#HEALTH INSURANCE CLAIM BLOCKCHAIN SOLUTION

### Use Case 1 - Out of Network Claims

* **Personas:** Insurance Payer
* **Technology:** Oracle Blockchain Platform, Oracle Health Insurance *(possibly)*
* **Outline:**

This demo will show how Oracle Blockchain Platform can be used to increase a payer's visibility into claims from out of network providers. We will take the example of a batch of claims from providers with which the payer has no managed care contract. On receiving the claims, the payer will write them to their blockchain network.

The blockchain record ids for the claims are sent to a third-party servicer that also has access to the blockchain network, and that servicer is also given access to write information to these claim records. That servicer will then begin their work on the claims, accessing and updating them on the blockchain network.

We'll show that, since the records remain on the blockchain network, the payer will always have the ability to check those records and see their current status. Once the servicer marks their work on the claims as complete, the payer will then initiate their payments.

If possible, this demo would show a version of Oracle Health Insurance that we've customized writing and reading the data to the blockchain network. If that is not possible, then we will build a similar front end and API to read and write the records from the blockchain network.

