/*
Package config provides capabilities to working with aws config
with sane configuration.

Welcome to the aws/config user guide. Here we will guide you through some
practical examples of how to interact with the API.

# Using aws client v2 start configuration

To provide flexibility in the API we offer some methods to allow additional
configuration during aws.config initialisation.

Use the api without any optional parameters.

		cfg, err := config.New()
	    if err != nil {
	        logger.Fatalf("unable to load aws SDK config, %v", err)
	    }
		// ... use cfg to initialize some aws client service

Using the optional API methods to configure region, profile and endpoint.

	cfg, err := config.New(
		config.WithEndpoint("http://mycustomendpoint:1234"),
		config.withRegion("us-west-2"),
	)
	if err != nil {
		logger.Fatalf("unable to load aws SDK config, %v", err)
	}
	// ... use cfg to initialize some aws client service

Note that you can use an optional parameter, e.g:

	cfg, err := config.New(
		config.withRegion("us-west-2"),
	)
	if err != nil {
		logger.Fatalf("unable to load aws SDK config, %v", err)
	}
	// ... use cfg to initialize some aws client service

# Using aws client v1 start configuration

You just have to change the name of the constructor to

	    sess, err := config.NewV1ClientConfig()
		if err != nil {
			logger.Fatalf("unable to load aws SDK config, %v", err)
		}
		// ... use sess to initialize some aws client service

Under the hood the API offers the same behaviour as the client for V2 of aws.
*/
package config
