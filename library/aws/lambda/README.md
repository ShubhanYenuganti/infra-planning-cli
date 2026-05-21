# Aws Lambda CLI

<fullname>Lambda</fullname> <p> <b>Overview</b> </p> <p>Lambda is a compute service that lets you run code without provisioning or managing servers. Lambda runs your code on a high-availability compute infrastructure and performs all of the administration of the compute resources, including server and operating system maintenance, capacity provisioning and automatic scaling, code monitoring and logging. With Lambda, you can run code for virtually any type of application or backend service. For more information about the Lambda service, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/welcome.html">What is Lambda</a> in the <b>Lambda Developer Guide</b>.</p> <p>The <i>Lambda API Reference</i> provides information about each of the API methods, including details about the parameters in each API request and response. </p> <p/> <p>You can use Software Development Kits (SDKs), Integrated Development Environment (IDE) Toolkits, and command line tools to access the API. For installation instructions, see <a href="http://aws.amazon.com/tools/">Tools for Amazon Web Services</a>. </p> <p>For a list of Region-specific endpoints that Lambda supports, see <a href="https://docs.aws.amazon.com/general/latest/gr/lambda-service.html/">Lambda endpoints and quotas </a> in the <i>Amazon Web Services General Reference.</i>. </p> <p>When making the API calls, you will need to authenticate your request by providing a signature. Lambda supports signature version 4. For more information, see <a href="https://docs.aws.amazon.com/general/latest/gr/signature-version-4.html">Signature Version 4 signing process</a> in the <i>Amazon Web Services General Reference.</i>. </p> <p> <b>CA certificates</b> </p> <p>Because Amazon Web Services SDKs use the CA certificates from your computer, changes to the certificates on the Amazon Web Services servers can cause connection failures when you attempt to use an SDK. You can prevent these failures by keeping your computer's CA certificates and operating system up-to-date. If you encounter this issue in a corporate environment and do not manage your own computer, you might need to ask an administrator to assist with the update process. The following list shows minimum operating system and Java versions:</p> <ul> <li> <p>Microsoft Windows versions that have updates from January 2005 or later installed contain at least one of the required CAs in their trust list. </p> </li> <li> <p>Mac OS X 10.4 with Java for Mac OS X 10.4 Release 5 (February 2007), Mac OS X 10.5 (October 2007), and later versions contain at least one of the required CAs in their trust list. </p> </li> <li> <p>Red Hat Enterprise Linux 5 (March 2007), 6, and 7 and CentOS 5, 6, and 7 all contain at least one of the required CAs in their default trusted CA list. </p> </li> <li> <p>Java 1.4.2_12 (May 2006), 5 Update 2 (March 2005), and all later versions, including Java 6 (December 2006), 7, and 8, contain at least one of the required CAs in their default trusted CA list. </p> </li> </ul> <p>When accessing the Lambda management console or Lambda API endpoints, whether through browsers or programmatically, you will need to ensure your client machines support any of the following CAs: </p> <ul> <li> <p>Amazon Root CA 1</p> </li> <li> <p>Starfield Services Root Certificate Authority - G2</p> </li> <li> <p>Starfield Class 2 Certification Authority</p> </li> </ul> <p>Root certificates from the first two authorities are available from <a href="https://www.amazontrust.com/repository/">Amazon trust services</a>, but keeping your computer up-to-date is the more straightforward solution. To learn more about ACM-provided certificates, see <a href="http://aws.amazon.com/certificate-manager/faqs/#certificates">Amazon Web Services Certificate Manager FAQs.</a> </p>

Learn more at [Aws Lambda](https://github.com/mermade/aws2openapi).

## Install

### Go

```
go install github.com/mvanhorn/printing-press-library/library/other/aws-lambda-pp-cli/cmd/aws-lambda-pp-cli@latest
```

### Binary

Download from [Releases](https://github.com/mvanhorn/printing-press-library/releases).

## Quick Start

### 1. Install

See [Install](#install) above.

### 2. Set Up Credentials

Get your API key from your API provider's developer portal. The key typically looks like a long alphanumeric string.

```bash
export AWS_LAMBDA_HMAC="<paste-your-key>"
```

You can also persist this in your config file at `~/.config/aws-lambda-pp-cli/config.toml`.

### 3. Verify Setup

```bash
aws-lambda-pp-cli doctor
```

This checks your configuration and credentials.

### 4. Try Your First Command

```bash
aws-lambda-pp-cli 2014-11-13 list
```

## Usage

<!-- HELP_OUTPUT -->

## Commands

### 2014-11-13

Manage 2014 11 13

- **`aws-lambda-pp-cli 2014-11-13 invoke-async`** - <important> <p>For asynchronous function invocation, use <a>Invoke</a>.</p> </important> <p>Invokes a function asynchronously.</p>

### 2015-03-31

Manage 2015 03 31

- **`aws-lambda-pp-cli 2015-03-31 add-permission`** - <p>Grants an Amazon Web Service, Amazon Web Services account, or Amazon Web Services organization permission to use a function. You can apply the policy at the function level, or specify a qualifier to restrict access to a single version or alias. If you use a qualifier, the invoker must use the full Amazon Resource Name (ARN) of that version or alias to invoke the function. Note: Lambda does not support adding policies to version $LATEST.</p> <p>To grant permission to another account, specify the account ID as the <code>Principal</code>. To grant permission to an organization defined in Organizations, specify the organization ID as the <code>PrincipalOrgID</code>. For Amazon Web Services, the principal is a domain-style identifier that the service defines, such as <code>s3.amazonaws.com</code> or <code>sns.amazonaws.com</code>. For Amazon Web Services, you can also specify the ARN of the associated resource as the <code>SourceArn</code>. If you grant permission to a service principal without specifying the source, other accounts could potentially configure resources in their account to invoke your Lambda function.</p> <p>This operation adds a statement to a resource-based permissions policy for the function. For more information about function policies, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/access-control-resource-based.html">Using resource-based policies for Lambda</a>.</p>
- **`aws-lambda-pp-cli 2015-03-31 create-alias`** - <p>Creates an <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-aliases.html">alias</a> for a Lambda function version. Use aliases to provide clients with a function identifier that you can update to invoke a different version.</p> <p>You can also map an alias to split invocation requests between two versions. Use the <code>RoutingConfig</code> parameter to specify a second version and the percentage of invocation requests that it receives.</p>
- **`aws-lambda-pp-cli 2015-03-31 create-event-source-mapping`** - <p>Creates a mapping between an event source and an Lambda function. Lambda reads items from the event source and invokes the function.</p> <p>For details about how to configure different event sources, see the following topics. </p> <ul> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-ddb.html#services-dynamodb-eventsourcemapping"> Amazon DynamoDB Streams</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-kinesis.html#services-kinesis-eventsourcemapping"> Amazon Kinesis</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html#events-sqs-eventsource"> Amazon SQS</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-mq.html#services-mq-eventsourcemapping"> Amazon MQ and RabbitMQ</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-msk.html"> Amazon MSK</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/kafka-smaa.html"> Apache Kafka</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-documentdb.html"> Amazon DocumentDB</a> </p> </li> </ul> <p>The following error handling options are available only for stream sources (DynamoDB and Kinesis):</p> <ul> <li> <p> <code>BisectBatchOnFunctionError</code> – If the function returns an error, split the batch in two and retry.</p> </li> <li> <p> <code>DestinationConfig</code> – Send discarded records to an Amazon SQS queue or Amazon SNS topic.</p> </li> <li> <p> <code>MaximumRecordAgeInSeconds</code> – Discard records older than the specified age. The default value is infinite (-1). When set to infinite (-1), failed records are retried until the record expires</p> </li> <li> <p> <code>MaximumRetryAttempts</code> – Discard records after the specified number of retries. The default value is infinite (-1). When set to infinite (-1), failed records are retried until the record expires.</p> </li> <li> <p> <code>ParallelizationFactor</code> – Process multiple batches from each shard concurrently.</p> </li> </ul> <p>For information about which configuration parameters apply to each event source, see the following topics.</p> <ul> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-ddb.html#services-ddb-params"> Amazon DynamoDB Streams</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-kinesis.html#services-kinesis-params"> Amazon Kinesis</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html#services-sqs-params"> Amazon SQS</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-mq.html#services-mq-params"> Amazon MQ and RabbitMQ</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-msk.html#services-msk-parms"> Amazon MSK</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-kafka.html#services-kafka-parms"> Apache Kafka</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-documentdb.html#docdb-configuration"> Amazon DocumentDB</a> </p> </li> </ul>
- **`aws-lambda-pp-cli 2015-03-31 create-function`** - <p>Creates a Lambda function. To create a function, you need a <a href="https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html">deployment package</a> and an <a href="https://docs.aws.amazon.com/lambda/latest/dg/intro-permission-model.html#lambda-intro-execution-role">execution role</a>. The deployment package is a .zip file archive or container image that contains your function code. The execution role grants the function permission to use Amazon Web Services, such as Amazon CloudWatch Logs for log streaming and X-Ray for request tracing.</p> <p>If the deployment package is a <a href="https://docs.aws.amazon.com/lambda/latest/dg/lambda-images.html">container image</a>, then you set the package type to <code>Image</code>. For a container image, the code property must include the URI of a container image in the Amazon ECR registry. You do not need to specify the handler and runtime properties.</p> <p>If the deployment package is a <a href="https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html#gettingstarted-package-zip">.zip file archive</a>, then you set the package type to <code>Zip</code>. For a .zip file archive, the code property specifies the location of the .zip file. You must also specify the handler and runtime properties. The code in the deployment package must be compatible with the target instruction set architecture of the function (<code>x86-64</code> or <code>arm64</code>). If you do not specify the architecture, then the default value is <code>x86-64</code>.</p> <p>When you create a function, Lambda provisions an instance of the function and its supporting resources. If your function connects to a VPC, this process can take a minute or so. During this time, you can't invoke or modify the function. The <code>State</code>, <code>StateReason</code>, and <code>StateReasonCode</code> fields in the response from <a>GetFunctionConfiguration</a> indicate when the function is ready to invoke. For more information, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/functions-states.html">Lambda function states</a>.</p> <p>A function has an unpublished version, and can have published versions and aliases. The unpublished version changes when you update your function's code and configuration. A published version is a snapshot of your function code and configuration that can't be changed. An alias is a named resource that maps to a version, and can be changed to map to a different version. Use the <code>Publish</code> parameter to create version <code>1</code> of your function from its initial configuration.</p> <p>The other parameters let you configure version-specific and function-level settings. You can modify version-specific settings later with <a>UpdateFunctionConfiguration</a>. Function-level settings apply to both the unpublished and published versions of the function, and include tags (<a>TagResource</a>) and per-function concurrency limits (<a>PutFunctionConcurrency</a>).</p> <p>You can use code signing if your deployment package is a .zip file archive. To enable code signing for this function, specify the ARN of a code-signing configuration. When a user attempts to deploy a code package with <a>UpdateFunctionCode</a>, Lambda checks that the code package has a valid signature from a trusted publisher. The code-signing configuration includes set of signing profiles, which define the trusted publishers for this function.</p> <p>If another Amazon Web Services account or an Amazon Web Service invokes your function, use <a>AddPermission</a> to grant permission by creating a resource-based Identity and Access Management (IAM) policy. You can grant permissions at the function level, on a version, or on an alias.</p> <p>To invoke your function directly, use <a>Invoke</a>. To invoke your function in response to events in other Amazon Web Services, create an event source mapping (<a>CreateEventSourceMapping</a>), or configure a function trigger in the other service. For more information, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/lambda-invocation.html">Invoking Lambda functions</a>.</p>
- **`aws-lambda-pp-cli 2015-03-31 delete-alias`** - Deletes a Lambda function <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-aliases.html">alias</a>.
- **`aws-lambda-pp-cli 2015-03-31 delete-event-source-mapping`** - <p>Deletes an <a href="https://docs.aws.amazon.com/lambda/latest/dg/intro-invocation-modes.html">event source mapping</a>. You can get the identifier of a mapping from the output of <a>ListEventSourceMappings</a>.</p> <p>When you delete an event source mapping, it enters a <code>Deleting</code> state and might not be completely deleted for several seconds.</p>
- **`aws-lambda-pp-cli 2015-03-31 delete-function`** - <p>Deletes a Lambda function. To delete a specific function version, use the <code>Qualifier</code> parameter. Otherwise, all versions and aliases are deleted.</p> <p>To delete Lambda event source mappings that invoke a function, use <a>DeleteEventSourceMapping</a>. For Amazon Web Services and resources that invoke your function directly, delete the trigger in the service where you originally configured it.</p>
- **`aws-lambda-pp-cli 2015-03-31 get-alias`** - Returns details about a Lambda function <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-aliases.html">alias</a>.
- **`aws-lambda-pp-cli 2015-03-31 get-event-source-mapping`** - Returns details about an event source mapping. You can get the identifier of a mapping from the output of <a>ListEventSourceMappings</a>.
- **`aws-lambda-pp-cli 2015-03-31 get-function`** - Returns information about the function or function version, with a link to download the deployment package that's valid for 10 minutes. If you specify a function version, only details that are specific to that version are returned.
- **`aws-lambda-pp-cli 2015-03-31 get-function-configuration`** - <p>Returns the version-specific settings of a Lambda function or version. The output includes only options that can vary between versions of a function. To modify these settings, use <a>UpdateFunctionConfiguration</a>.</p> <p>To get all of a function's details, including function-level settings, use <a>GetFunction</a>.</p>
- **`aws-lambda-pp-cli 2015-03-31 get-policy`** - Returns the <a href="https://docs.aws.amazon.com/lambda/latest/dg/access-control-resource-based.html">resource-based IAM policy</a> for a function, version, or alias.
- **`aws-lambda-pp-cli 2015-03-31 invoke`** - <p>Invokes a Lambda function. You can invoke a function synchronously (and wait for the response), or asynchronously. To invoke a function asynchronously, set <code>InvocationType</code> to <code>Event</code>.</p> <p>For <a href="https://docs.aws.amazon.com/lambda/latest/dg/invocation-sync.html">synchronous invocation</a>, details about the function response, including errors, are included in the response body and headers. For either invocation type, you can find more information in the <a href="https://docs.aws.amazon.com/lambda/latest/dg/monitoring-functions.html">execution log</a> and <a href="https://docs.aws.amazon.com/lambda/latest/dg/lambda-x-ray.html">trace</a>.</p> <p>When an error occurs, your function may be invoked multiple times. Retry behavior varies by error type, client, event source, and invocation type. For example, if you invoke a function asynchronously and it returns an error, Lambda executes the function up to two more times. For more information, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/invocation-retries.html">Error handling and automatic retries in Lambda</a>.</p> <p>For <a href="https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html">asynchronous invocation</a>, Lambda adds events to a queue before sending them to your function. If your function does not have enough capacity to keep up with the queue, events may be lost. Occasionally, your function may receive the same event multiple times, even if no error occurs. To retain events that were not processed, configure your function with a <a href="https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html#invocation-dlq">dead-letter queue</a>.</p> <p>The status code in the API response doesn't reflect function errors. Error codes are reserved for errors that prevent your function from executing, such as permissions errors, <a href="https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-limits.html">quota</a> errors, or issues with your function's code and configuration. For example, Lambda returns <code>TooManyRequestsException</code> if running the function would cause you to exceed a concurrency limit at either the account level (<code>ConcurrentInvocationLimitExceeded</code>) or function level (<code>ReservedFunctionConcurrentInvocationLimitExceeded</code>).</p> <p>For functions with a long timeout, your client might disconnect during synchronous invocation while it waits for a response. Configure your HTTP client, SDK, firewall, proxy, or operating system to allow for long connections with timeout or keep-alive settings.</p> <p>This operation requires permission for the <a href="https://docs.aws.amazon.com/IAM/latest/UserGuide/list_awslambda.html">lambda:InvokeFunction</a> action. For details on how to set up permissions for cross-account invocations, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/access-control-resource-based.html#permissions-resource-xaccountinvoke">Granting function access to other accounts</a>.</p>
- **`aws-lambda-pp-cli 2015-03-31 list-aliases`** - Returns a list of <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-aliases.html">aliases</a> for a Lambda function.
- **`aws-lambda-pp-cli 2015-03-31 list-event-source-mappings`** - Lists event source mappings. Specify an <code>EventSourceArn</code> to show only event source mappings for a single event source.
- **`aws-lambda-pp-cli 2015-03-31 list-functions`** - <p>Returns a list of Lambda functions, with the version-specific configuration of each. Lambda returns up to 50 functions per call.</p> <p>Set <code>FunctionVersion</code> to <code>ALL</code> to include all published versions of each function in addition to the unpublished version.</p> <note> <p>The <code>ListFunctions</code> operation returns a subset of the <a>FunctionConfiguration</a> fields. To get the additional fields (State, StateReasonCode, StateReason, LastUpdateStatus, LastUpdateStatusReason, LastUpdateStatusReasonCode, RuntimeVersionConfig) for a function or version, use <a>GetFunction</a>.</p> </note>
- **`aws-lambda-pp-cli 2015-03-31 list-versions-by-function`** - Returns a list of <a href="https://docs.aws.amazon.com/lambda/latest/dg/versioning-aliases.html">versions</a>, with the version-specific configuration of each. Lambda returns up to 50 versions per call.
- **`aws-lambda-pp-cli 2015-03-31 publish-version`** - <p>Creates a <a href="https://docs.aws.amazon.com/lambda/latest/dg/versioning-aliases.html">version</a> from the current code and configuration of a function. Use versions to create a snapshot of your function code and configuration that doesn't change.</p> <p>Lambda doesn't publish a version if the function's configuration and code haven't changed since the last version. Use <a>UpdateFunctionCode</a> or <a>UpdateFunctionConfiguration</a> to update the function before publishing a version.</p> <p>Clients can invoke versions directly or with an alias. To create an alias, use <a>CreateAlias</a>.</p>
- **`aws-lambda-pp-cli 2015-03-31 remove-permission`** - Revokes function-use permission from an Amazon Web Service or another Amazon Web Services account. You can get the ID of the statement from the output of <a>GetPolicy</a>.
- **`aws-lambda-pp-cli 2015-03-31 update-alias`** - Updates the configuration of a Lambda function <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-aliases.html">alias</a>.
- **`aws-lambda-pp-cli 2015-03-31 update-event-source-mapping`** - <p>Updates an event source mapping. You can change the function that Lambda invokes, or pause invocation and resume later from the same location.</p> <p>For details about how to configure different event sources, see the following topics. </p> <ul> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-ddb.html#services-dynamodb-eventsourcemapping"> Amazon DynamoDB Streams</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-kinesis.html#services-kinesis-eventsourcemapping"> Amazon Kinesis</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html#events-sqs-eventsource"> Amazon SQS</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-mq.html#services-mq-eventsourcemapping"> Amazon MQ and RabbitMQ</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-msk.html"> Amazon MSK</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/kafka-smaa.html"> Apache Kafka</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-documentdb.html"> Amazon DocumentDB</a> </p> </li> </ul> <p>The following error handling options are available only for stream sources (DynamoDB and Kinesis):</p> <ul> <li> <p> <code>BisectBatchOnFunctionError</code> – If the function returns an error, split the batch in two and retry.</p> </li> <li> <p> <code>DestinationConfig</code> – Send discarded records to an Amazon SQS queue or Amazon SNS topic.</p> </li> <li> <p> <code>MaximumRecordAgeInSeconds</code> – Discard records older than the specified age. The default value is infinite (-1). When set to infinite (-1), failed records are retried until the record expires</p> </li> <li> <p> <code>MaximumRetryAttempts</code> – Discard records after the specified number of retries. The default value is infinite (-1). When set to infinite (-1), failed records are retried until the record expires.</p> </li> <li> <p> <code>ParallelizationFactor</code> – Process multiple batches from each shard concurrently.</p> </li> </ul> <p>For information about which configuration parameters apply to each event source, see the following topics.</p> <ul> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-ddb.html#services-ddb-params"> Amazon DynamoDB Streams</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-kinesis.html#services-kinesis-params"> Amazon Kinesis</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-sqs.html#services-sqs-params"> Amazon SQS</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-mq.html#services-mq-params"> Amazon MQ and RabbitMQ</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-msk.html#services-msk-parms"> Amazon MSK</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-kafka.html#services-kafka-parms"> Apache Kafka</a> </p> </li> <li> <p> <a href="https://docs.aws.amazon.com/lambda/latest/dg/with-documentdb.html#docdb-configuration"> Amazon DocumentDB</a> </p> </li> </ul>
- **`aws-lambda-pp-cli 2015-03-31 update-function-code`** - <p>Updates a Lambda function's code. If code signing is enabled for the function, the code package must be signed by a trusted publisher. For more information, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-codesigning.html">Configuring code signing for Lambda</a>.</p> <p>If the function's package type is <code>Image</code>, then you must specify the code package in <code>ImageUri</code> as the URI of a <a href="https://docs.aws.amazon.com/lambda/latest/dg/lambda-images.html">container image</a> in the Amazon ECR registry.</p> <p>If the function's package type is <code>Zip</code>, then you must specify the deployment package as a <a href="https://docs.aws.amazon.com/lambda/latest/dg/gettingstarted-package.html#gettingstarted-package-zip">.zip file archive</a>. Enter the Amazon S3 bucket and key of the code .zip file location. You can also provide the function code inline using the <code>ZipFile</code> field.</p> <p>The code in the deployment package must be compatible with the target instruction set architecture of the function (<code>x86-64</code> or <code>arm64</code>).</p> <p>The function's code is locked when you publish a version. You can't modify the code of a published version, only the unpublished version.</p> <note> <p>For a function defined as a container image, Lambda resolves the image tag to an image digest. In Amazon ECR, if you update the image tag to a new image, Lambda does not automatically update the function.</p> </note>
- **`aws-lambda-pp-cli 2015-03-31 update-function-configuration`** - <p>Modify the version-specific settings of a Lambda function.</p> <p>When you update a function, Lambda provisions an instance of the function and its supporting resources. If your function connects to a VPC, this process can take a minute. During this time, you can't modify the function, but you can still invoke it. The <code>LastUpdateStatus</code>, <code>LastUpdateStatusReason</code>, and <code>LastUpdateStatusReasonCode</code> fields in the response from <a>GetFunctionConfiguration</a> indicate when the update is complete and the function is processing events with the new configuration. For more information, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/functions-states.html">Lambda function states</a>.</p> <p>These settings can vary between versions of a function and are locked when you publish a version. You can't modify the configuration of a published version, only the unpublished version.</p> <p>To configure function concurrency, use <a>PutFunctionConcurrency</a>. To grant invoke permissions to an Amazon Web Services account or Amazon Web Service, use <a>AddPermission</a>.</p>

### 2016-08-19

Manage 2016 08 19

- **`aws-lambda-pp-cli 2016-08-19 get-account-settings`** - Retrieves details about your account's <a href="https://docs.aws.amazon.com/lambda/latest/dg/limits.html">limits</a> and usage in an Amazon Web Services Region.

### 2017-03-31

Manage 2017 03 31

- **`aws-lambda-pp-cli 2017-03-31 list-tags`** - Returns a function's <a href="https://docs.aws.amazon.com/lambda/latest/dg/tagging.html">tags</a>. You can also view tags with <a>GetFunction</a>.
- **`aws-lambda-pp-cli 2017-03-31 tag-resource`** - Adds <a href="https://docs.aws.amazon.com/lambda/latest/dg/tagging.html">tags</a> to a function.
- **`aws-lambda-pp-cli 2017-03-31 untag-resource`** - Removes <a href="https://docs.aws.amazon.com/lambda/latest/dg/tagging.html">tags</a> from a function.

### 2017-10-31

Manage 2017 10 31

- **`aws-lambda-pp-cli 2017-10-31 delete-function-concurrency`** - Removes a concurrent execution limit from a function.
- **`aws-lambda-pp-cli 2017-10-31 put-function-concurrency`** - <p>Sets the maximum number of simultaneous executions for a function, and reserves capacity for that concurrency level.</p> <p>Concurrency settings apply to the function as a whole, including all published versions and the unpublished version. Reserving concurrency both ensures that your function has capacity to process the specified number of events simultaneously, and prevents it from scaling beyond that level. Use <a>GetFunction</a> to see the current setting for a function.</p> <p>Use <a>GetAccountSettings</a> to see your Regional concurrency limit. You can reserve concurrency for as many functions as you like, as long as you leave at least 100 simultaneous executions unreserved for functions that aren't configured with a per-function limit. For more information, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/invocation-scaling.html">Lambda function scaling</a>.</p>

### 2018-10-31

Manage 2018 10 31

- **`aws-lambda-pp-cli 2018-10-31 add-layer-version-permission`** - <p>Adds permissions to the resource-based policy of a version of an <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html">Lambda layer</a>. Use this action to grant layer usage permission to other accounts. You can grant permission to a single account, all accounts in an organization, or all Amazon Web Services accounts. </p> <p>To revoke permission, call <a>RemoveLayerVersionPermission</a> with the statement ID that you specified when you added it.</p>
- **`aws-lambda-pp-cli 2018-10-31 delete-layer-version`** - Deletes a version of an <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html">Lambda layer</a>. Deleted versions can no longer be viewed or added to functions. To avoid breaking functions, a copy of the version remains in Lambda until no functions refer to it.
- **`aws-lambda-pp-cli 2018-10-31 get-layer-version`** - Returns information about a version of an <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html">Lambda layer</a>, with a link to download the layer archive that's valid for 10 minutes.
- **`aws-lambda-pp-cli 2018-10-31 get-layer-version-by-arn`** - Returns information about a version of an <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html">Lambda layer</a>, with a link to download the layer archive that's valid for 10 minutes.
- **`aws-lambda-pp-cli 2018-10-31 get-layer-version-policy`** - Returns the permission policy for a version of an <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html">Lambda layer</a>. For more information, see <a>AddLayerVersionPermission</a>.
- **`aws-lambda-pp-cli 2018-10-31 list-layer-versions`** - Lists the versions of an <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html">Lambda layer</a>. Versions that have been deleted aren't listed. Specify a <a href="https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html">runtime identifier</a> to list only versions that indicate that they're compatible with that runtime. Specify a compatible architecture to include only layer versions that are compatible with that architecture.
- **`aws-lambda-pp-cli 2018-10-31 list-layers`** - Lists <a href="https://docs.aws.amazon.com/lambda/latest/dg/invocation-layers.html">Lambda layers</a> and shows information about the latest version of each. Specify a <a href="https://docs.aws.amazon.com/lambda/latest/dg/lambda-runtimes.html">runtime identifier</a> to list only layers that indicate that they're compatible with that runtime. Specify a compatible architecture to include only layers that are compatible with that <a href="https://docs.aws.amazon.com/lambda/latest/dg/foundation-arch.html">instruction set architecture</a>.
- **`aws-lambda-pp-cli 2018-10-31 publish-layer-version`** - <p>Creates an <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html">Lambda layer</a> from a ZIP archive. Each time you call <code>PublishLayerVersion</code> with the same layer name, a new version is created.</p> <p>Add layers to your function with <a>CreateFunction</a> or <a>UpdateFunctionConfiguration</a>.</p>
- **`aws-lambda-pp-cli 2018-10-31 remove-layer-version-permission`** - Removes a statement from the permissions policy for a version of an <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-layers.html">Lambda layer</a>. For more information, see <a>AddLayerVersionPermission</a>.

### 2019-09-25

Manage 2019 09 25

- **`aws-lambda-pp-cli 2019-09-25 delete-function-event-invoke-config`** - <p>Deletes the configuration for asynchronous invocation for a function, version, or alias.</p> <p>To configure options for asynchronous invocation, use <a>PutFunctionEventInvokeConfig</a>.</p>
- **`aws-lambda-pp-cli 2019-09-25 get-function-event-invoke-config`** - <p>Retrieves the configuration for asynchronous invocation for a function, version, or alias.</p> <p>To configure options for asynchronous invocation, use <a>PutFunctionEventInvokeConfig</a>.</p>
- **`aws-lambda-pp-cli 2019-09-25 list-function-event-invoke-configs`** - <p>Retrieves a list of configurations for asynchronous invocation for a function.</p> <p>To configure options for asynchronous invocation, use <a>PutFunctionEventInvokeConfig</a>.</p>
- **`aws-lambda-pp-cli 2019-09-25 put-function-event-invoke-config`** - <p>Configures options for <a href="https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html">asynchronous invocation</a> on a function, version, or alias. If a configuration already exists for a function, version, or alias, this operation overwrites it. If you exclude any settings, they are removed. To set one option without affecting existing settings for other options, use <a>UpdateFunctionEventInvokeConfig</a>.</p> <p>By default, Lambda retries an asynchronous invocation twice if the function returns an error. It retains events in a queue for up to six hours. When an event fails all processing attempts or stays in the asynchronous invocation queue for too long, Lambda discards it. To retain discarded events, configure a dead-letter queue with <a>UpdateFunctionConfiguration</a>.</p> <p>To send an invocation record to a queue, topic, function, or event bus, specify a <a href="https://docs.aws.amazon.com/lambda/latest/dg/invocation-async.html#invocation-async-destinations">destination</a>. You can configure separate destinations for successful invocations (on-success) and events that fail all processing attempts (on-failure). You can configure destinations in addition to or instead of a dead-letter queue.</p>
- **`aws-lambda-pp-cli 2019-09-25 update-function-event-invoke-config`** - <p>Updates the configuration for asynchronous invocation for a function, version, or alias.</p> <p>To configure options for asynchronous invocation, use <a>PutFunctionEventInvokeConfig</a>.</p>

### 2019-09-30

Manage 2019 09 30

- **`aws-lambda-pp-cli 2019-09-30 delete-provisioned-concurrency-config`** - Deletes the provisioned concurrency configuration for a function.
- **`aws-lambda-pp-cli 2019-09-30 get-function-concurrency`** - Returns details about the reserved concurrency configuration for a function. To set a concurrency limit for a function, use <a>PutFunctionConcurrency</a>.
- **`aws-lambda-pp-cli 2019-09-30 get-provisioned-concurrency-config`** - Retrieves the provisioned concurrency configuration for a function's alias or version.
- **`aws-lambda-pp-cli 2019-09-30 list-provisioned-concurrency-configs`** - Retrieves a list of provisioned concurrency configurations for a function.
- **`aws-lambda-pp-cli 2019-09-30 put-provisioned-concurrency-config`** - Adds a provisioned concurrency configuration to a function's alias or version.

### 2020-04-22

Manage 2020 04 22

- **`aws-lambda-pp-cli 2020-04-22 create-code-signing-config`** - Creates a code signing configuration. A <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-codesigning.html">code signing configuration</a> defines a list of allowed signing profiles and defines the code-signing validation policy (action to be taken if deployment validation checks fail).
- **`aws-lambda-pp-cli 2020-04-22 delete-code-signing-config`** - Deletes the code signing configuration. You can delete the code signing configuration only if no function is using it.
- **`aws-lambda-pp-cli 2020-04-22 get-code-signing-config`** - Returns information about the specified code signing configuration.
- **`aws-lambda-pp-cli 2020-04-22 list-code-signing-configs`** - Returns a list of <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuring-codesigning.html">code signing configurations</a>. A request returns up to 10,000 configurations per call. You can use the <code>MaxItems</code> parameter to return fewer configurations per call.
- **`aws-lambda-pp-cli 2020-04-22 list-functions-by-code-signing-config`** - List the functions that use the specified code signing configuration. You can use this method prior to deleting a code signing configuration, to verify that no functions are using it.
- **`aws-lambda-pp-cli 2020-04-22 update-code-signing-config`** - Update the code signing configuration. Changes to the code signing configuration take effect the next time a user tries to deploy a code package to the function.

### 2020-06-30

Manage 2020 06 30

- **`aws-lambda-pp-cli 2020-06-30 delete-function-code-signing-config`** - Removes the code signing configuration from the function.
- **`aws-lambda-pp-cli 2020-06-30 get-function-code-signing-config`** - Returns the code signing configuration for the specified function.
- **`aws-lambda-pp-cli 2020-06-30 put-function-code-signing-config`** - Update the code signing configuration for the function. Changes to the code signing configuration take effect the next time a user tries to deploy a code package to the function.

### 2021-07-20

Manage 2021 07 20

- **`aws-lambda-pp-cli 2021-07-20 get-runtime-management-config`** - Retrieves the runtime management configuration for a function's version. If the runtime update mode is <b>Manual</b>, this includes the ARN of the runtime version and the runtime update mode. If the runtime update mode is <b>Auto</b> or <b>Function update</b>, this includes the runtime update mode and <code>null</code> is returned for the ARN. For more information, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/runtimes-update.html">Runtime updates</a>.
- **`aws-lambda-pp-cli 2021-07-20 put-runtime-management-config`** - Sets the runtime management configuration for a function's version. For more information, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/runtimes-update.html">Runtime updates</a>.

### 2021-10-31

Manage 2021 10 31

- **`aws-lambda-pp-cli 2021-10-31 create-function-url-config`** - Creates a Lambda function URL with the specified configuration parameters. A function URL is a dedicated HTTP(S) endpoint that you can use to invoke your function.
- **`aws-lambda-pp-cli 2021-10-31 delete-function-url-config`** - Deletes a Lambda function URL. When you delete a function URL, you can't recover it. Creating a new function URL results in a different URL address.
- **`aws-lambda-pp-cli 2021-10-31 get-function-url-config`** - Returns details about a Lambda function URL.
- **`aws-lambda-pp-cli 2021-10-31 list-function-url-configs`** - Returns a list of Lambda function URLs for the specified function.
- **`aws-lambda-pp-cli 2021-10-31 update-function-url-config`** - Updates the configuration for a Lambda function URL.

### 2021-11-15

Manage 2021 11 15

- **`aws-lambda-pp-cli 2021-11-15 invoke-with-response-stream`** - <p>Configure your Lambda functions to stream response payloads back to clients. For more information, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/configuration-response-streaming.html">Configuring a Lambda function to stream responses</a>.</p> <p>This operation requires permission for the <a href="https://docs.aws.amazon.com/IAM/latest/UserGuide/list_awslambda.html">lambda:InvokeFunction</a> action. For details on how to set up permissions for cross-account invocations, see <a href="https://docs.aws.amazon.com/lambda/latest/dg/access-control-resource-based.html#permissions-resource-xaccountinvoke">Granting function access to other accounts</a>.</p>


## Output Formats

```bash
# Human-readable table (default in terminal, JSON when piped)
aws-lambda-pp-cli 2014-11-13 list

# JSON for scripting and agents
aws-lambda-pp-cli 2014-11-13 list --json

# Filter to specific fields
aws-lambda-pp-cli 2014-11-13 list --json --select id,name,status

# Dry run — show the request without sending
aws-lambda-pp-cli 2014-11-13 list --dry-run

# Agent mode — JSON + compact + no prompts in one flag
aws-lambda-pp-cli 2014-11-13 list --agent
```

## Agent Usage

This CLI is designed for AI agent consumption:

- **Non-interactive** - never prompts, every input is a flag
- **Pipeable** - `--json` output to stdout, errors to stderr
- **Filterable** - `--select id,name` returns only fields you need
- **Previewable** - `--dry-run` shows the request without sending
- **Retryable** - creates return "already exists" on retry, deletes return "already deleted"
- **Confirmable** - `--yes` for explicit confirmation of destructive actions
- **Piped input** - `echo '{"key":"value"}' | aws-lambda-pp-cli <resource> create --stdin`
- **Cacheable** - GET responses cached for 5 minutes, bypass with `--no-cache`
- **Agent-safe by default** - no colors or formatting unless `--human-friendly` is set
- **Progress events** - paginated commands emit NDJSON events to stderr in default mode

Exit codes: `0` success, `2` usage error, `3` not found, `4` auth error, `5` API error, `7` rate limited, `10` config error.

## Use as MCP Server

This CLI ships a companion MCP server for use with Claude Desktop, Cursor, and other MCP-compatible tools.

### Claude Code

```bash
claude mcp add aws-lambda aws-lambda-pp-mcp -e AWS_LAMBDA_HMAC=<your-key>
```

### Claude Desktop

Add to your Claude Desktop config (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "aws-lambda": {
      "command": "aws-lambda-pp-mcp",
      "env": {
        "AWS_LAMBDA_HMAC": "<your-key>"
      }
    }
  }
}
```

## Cookbook

Common workflows and recipes:

```bash
# List resources as JSON for scripting
aws-lambda-pp-cli 2014-11-13 list --json

# Filter to specific fields
aws-lambda-pp-cli 2014-11-13 list --json --select id,name,status

# Dry run to preview the request
aws-lambda-pp-cli 2014-11-13 list --dry-run

# Sync data locally for offline search
aws-lambda-pp-cli sync

# Search synced data
aws-lambda-pp-cli search "query"

# Export for backup
aws-lambda-pp-cli export --format jsonl > backup.jsonl
```

## Health Check

```bash
aws-lambda-pp-cli doctor
```

<!-- DOCTOR_OUTPUT -->

## Configuration

Config file: `~/.config/aws-lambda-pp-cli/config.toml`

Environment variables:
- `AWS_LAMBDA_HMAC`

## Troubleshooting

**Authentication errors (exit code 4)**
- Run `aws-lambda-pp-cli doctor` to check credentials
- Verify the environment variable is set: `echo $AWS_LAMBDA_HMAC`

**Not found errors (exit code 3)**
- Check the resource ID is correct
- Run the `list` command to see available items

**Rate limit errors (exit code 7)**
- The CLI auto-retries with exponential backoff
- If persistent, wait a few minutes and try again

---

Generated by [CLI Printing Press](https://github.com/mvanhorn/cli-printing-press)
