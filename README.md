# Abios Push API Controller
Utility program to manage Abios push API subscriptions, as well as other push API related commands.

## Requirements
You need to have valid Abios API keys to run this application. If you don't have any keys, please contact us at `info@abiosgaming.com` and we'll help you to get setup.
 
This application has been tested with Golang 1.9.x, it might work with older compiler versions.

All external library dependencies are included in the `vendor` directory. If you need to reinstall them for some reason, remove the `vendor` directory and regenerate it using the `glide` dependency management tool (see `https://glide.sh` for info on how to install it).

## Compiling
To compile the client:

`$ go build .`

Now you should have a binary called `push-api-ctl`.


If you want to reinstall the library dependencies, do:

`$ glide install`

This creates the `vendor` directory with all the dependencies.


## Usage

Below is a list of all commands available. The help text for each command (e.g. `push-api-ctl create --help`) shows more information
about option flags.


 * Show push api account configuration
    `./push-api-ctl --client-id=<...> --client-secret=<...> config`

 * List all currently registered subscriptions

    `./push-api-ctl --client-id=<...> --client-secret=<...> list`

 * Create a new subscription and register it with the Push Service
 
    `./push-api-ctl --client-id=<...> --client-secret=<...> create -f <Subscription JSON spec file>`
    The command takes a subscription specification in JSON format as input, see [the subscription documentation](https://docs.abiosgaming.com/v2/reference#section-2-subscription-specifications) for information about how to write subscription specifications.

    You can optionally use the `--name=<Name>` option to give the subscription a name.
 
 * Show subscription specification for a registered subscription

    `./push-api-ctl --client-id=<...> --client-secret=<...> get <NameOrID>`
    The `<NameOrID>` argument is either the UUID of the subscription or the name used when creating/registering the subscription.

 * Update subscription

    `./push-api-ctl --client-id=<...> --client-secret=<...> update -f <Subscription JSON spec file> <NameOrID>`
    Updates an existing subscription with a new specification.

 * Delete subscription

    `./push-api-ctl --client-id=<...> --client-secret=<...> <NameOrID>...`
    
    Delete one more more subscriptions with given id's or names.
