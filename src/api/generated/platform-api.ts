// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2026 Steadybit GmbH

// Generated from openapi/platform-api.json by `npm run api:generate`. Do not edit.

export interface paths {
    "/api/access-tokens": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get access token list
         * @deprecated
         * @description Deprecated, use v2 instead. Get a list of all access tokens. The access token itself is abbreviated for security reasons. Access tokens with v2 features are not returned, as they can not be represented cleanly in the old format.
         */
        get: operations["getAccessTokens"];
        put?: never;
        /**
         * Add a access token
         * @deprecated
         * @description Deprecated, use v2 instead. Generate a new access token associated to. This access token can be used for e.g. creating new experiments and running experiments.
         */
        post: operations["createAccessToken"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/access-tokens/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /**
         * Delete access token
         * @deprecated
         * @description Remove the access token associated. After that, the access token can't be used anymore for e.g. creating a new experiment or running an experiment.
         */
        delete: operations["deleteAccessToken"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/access-tokens/v2": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get access token list
         * @description Get a list of all access tokens. The access token itself is abbreviated for security reasons.
         */
        get: operations["getAccessTokens_1"];
        put?: never;
        /**
         * Create an access token
         * @description Generate a new access token.
         */
        post: operations["createAccessToken_1"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/access-tokens/v2/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /**
         * Delete access token
         * @description Remove the access token. After that, the access token can't be used anymore.
         */
        delete: operations["deleteAccessToken_1"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/access-tokens/v2/{id}/recreate": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Recreate an access token
         * @description Recreate an existing access token with a new expiration date. The old token is deleted and a new one is generated with the same name, type, and team associations.
         */
        post: operations["recreateAccessToken"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/actions": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get all actions. */
        get: operations["findAllActions"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/actions/{actionId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single action description
         * @description Get action including their parameters.
         */
        get: operations["getAction"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/advice": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Get all currently active advice for a given environment and query. */
        post: operations["getTargetAdviceSummary"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/audit-log": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get all audit log entries
         * @description Retrieve all audit logs in the given time-frame.<br/>This endpoint requires an admin-token and can't be used with a team-based token.
         */
        get: operations["find"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/badges/link": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Forward to Steadybit platform to either create an experiment associated to the `tag` or forward to the experiments linked already to the `tag`
         * @description This endpoint can be used as a link for the badge of the `/api/badges/linked-badge.svg` API to either create a new experiment or show the linked experiments in Steadybit. This will help to link it correctly e.g. in your CMS-systems.
         */
        get: operations["forwardToPlatform"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/badges/linked-badge.svg": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get badge for create experiment or run status as SVG image
         * @description Creates an image badge that is either for creating a new experiment linked to an `externalReference` or - if an experiment with the given `externalReference` already exists - a badge showing the run status of the experiment. The badge is return as SVG to integrate it nicely e.g. into your CMS-systems. You can use the `/api/badges/link` endpoint to link it appropriately
         */
        get: operations["getLinkedBadge"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/environments": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of all environments
         * @description Get a list of all environments that exist.
         */
        get: operations["getEnvironments"];
        put?: never;
        /**
         * Create or update an environment
         * @description Insert or update the environment in Steadybit. The `id` will be used to identify whether the environment exists already and should be updated or newly inserted.
         */
        post: operations["upsertEnvironment"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/environments/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single environment
         * @description Get all details of a single existing environment.
         */
        get: operations["getEnvironment"];
        put?: never;
        post?: never;
        /**
         * Delete environment
         * @description Remove the given environment from the Steadybit platform.
         */
        delete: operations["deleteEnvironment"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/environments/{id}/variables": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get environment variables
         * @description Get all environment variables associated to a single environment
         */
        get: operations["getEnvironmentVariables"];
        /**
         * Add / merge all environment variables
         * @description All provided environment variables will be associated to the given environment.<br/>If an environment variable key is already in use, it's value is updated with the value being provided.<br/>If an environment variable is already associated to the environment but not provided, it continues to exist.
         */
        put: operations["updateEnvironmentVariables"];
        /**
         * Replace all environment variables
         * @description All provided environment variables will be associated to the given environment and existing ones removed.<br/>If an environment variable key is already in use, it's value is updated with the value being provided.<br/>If an environment variable is already associated to the environment but not provided, it will be removed.
         */
        post: operations["setEnvironmentVariables"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of all experiments
         * @description Get a list of all experiments that exist.
         */
        get: operations["getExperiments"];
        put?: never;
        /**
         * Create or update an experiment
         * @description Insert or update the experiment. The `externalId` will be used to identify whether the experiment exists already and should be updated or newly inserted.
         */
        post: operations["createOrUpdateExperiment"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/{key}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single experiment
         * @description Get all details of a single existing experiment.
         */
        get: operations["getExperiment"];
        put?: never;
        /**
         * Update an experiment
         * @description Update the experiment identified by the experiment `key`.
         */
        post: operations["updateExperiment"];
        /**
         * Delete experiment
         * @description Remove the given experiment. The associated number is still reserved afterwards and will not be reused.
         */
        delete: operations["deleteExperiment"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/{key}/badge.svg": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get experiment run status as SVG image
         * @description Get the status of the latest experiment run of the associated experiment as SVG to integrate it nicely e.g. into your CMS-systems.
         */
        get: operations["getExperimentBadge"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/{key}/execute": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Execute an experiment
         * @description Trigger execution of a single experiment specified by `key`. The body is optional and allows to specify overrides and custom properties for the experiment execution.
         *
         *        Examples:
         *        - Override environment from the experiment for a single run:
         *        ```
         *        {
         *            "environment":  "Shop Stage"
         *        }
         *        ```
         *        - Override the variables for a single execution:
         *        ```
         *        {
         *            "variables": {
         *                "httpEndpoint": "http://dev.shop.products.internal"
         *            }
         *        }
         *        ```
         */
        post: operations["executeExperiment"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/{key}/executions": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of all experiment executions of a single experiment
         * @description Get a list of all experiment executions that were performed for a specific experiment.
         */
        get: operations["getExperimentExecutions_3"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/execute": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Save and run experiment
         * @description Save the given experiment and immediately run it.
         */
        post: operations["saveAndRun"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/executions": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of all experiment executions
         * @deprecated
         * @description Get a list of all experiment executions that exist.
         */
        get: operations["getExperimentExecutions_1"];
        put?: never;
        /**
         * Fetch a list of experiment executions
         * @description Get list of experiment executions given a set of filters. The result is sorted by creation date in descending order. The result is paged with a page size of 50.
         */
        post: operations["getExperimentExecutions_2"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/executions/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single experiment executions of a single experiment
         * @description Get a single experiment execution that was performed for a specific experiment.
         */
        get: operations["getExperimentExecution"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/executions/{id}/artifacts/{targetExecutionId}/{artifactId}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get: operations["getArtifact"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/executions/{id}/cancel": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Cancel a running experiment execution of a single experiment
         * @description Cancels a currently running experiment execution to be stopped as soon as possible.
         */
        post: operations["cancelExperimentExecution"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/executions/{id}/properties": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Update properties of an experiment execution
         * @description Update properties of an experiment execution. This is only possible for associated properties with `editableInExecution` set to `true` or for properties that have been added after the execution.
         */
        post: operations["updateExecutionProperties"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/executions/{id}/properties/{key}/add": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Add a single value to a list property of an experiment execution.
         * @description This operation will fail if the property identified by `key` is not a list property. Only properties with `editableInExecution` set to `true` can be modified.
         */
        post: operations["addExecutionPropertyValue"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/executions/{id}/properties/{key}/set": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Set the value of a property of an experiment execution.
         * @description Only properties with `editableInExecution` set to `true` can be modified.
         */
        post: operations["setExecutionPropertyValue"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/schedules": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Create or update an experiment schedule
         * @description Insert or update the experiment schedule. The `id` will be used to identify whether the schedule exists already and should be updated or newly inserted.
         */
        post: operations["upsertSchedule"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/schedules/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get experiment schedules for a specific experiment schedule id */
        get: operations["getSchedules"];
        put?: never;
        post?: never;
        /** Remove an existing experiment schedule */
        delete: operations["removeExperimentScheduleById"];
        options?: never;
        head?: never;
        /**
         * Partially update an experiment schedule
         * @description Update specific fields of an existing experiment schedule. Only non-null fields in the request body will be updated.
         */
        patch: operations["patchSchedule"];
        trace?: never;
    };
    "/api/experiments/schedules/v2": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get all current experiment schedule configurations */
        get: operations["getAllSchedulesV2"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/templates": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of all templates
         * @description Get a list of all templates that exist.
         */
        get: operations["getExperimentTemplates"];
        put?: never;
        /**
         * Create or update an experiment template
         * @description Insert or update the experiment template in Steadybit. The `id` will be used to identify whether the experiment template exists already and should be updated or newly inserted. If this template is used in a service profile, existing provided service experiments will get updated.
         */
        post: operations["upsertExperimentTemplate"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/templates/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single experiment template
         * @description Get all details of a single existing experiment template.
         */
        get: operations["getExperimentTemplate"];
        put?: never;
        post?: never;
        /**
         * Delete experiment template
         * @description Remove the given experiment template from the Steadybit platform. If this template is used in a service profile, it will be removed from the profile and all provided service experiments will get deleted.
         */
        delete: operations["deleteExperimentTemplate"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/templates/{id}/experiment-create": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Create an experiment based on an experiment template
         * @description Use the given experiment template id and the placeholder values to create or update the experiment. The `externalId` will be used to identify whether the experiment exists already and should be updated or newly inserted.
         */
        post: operations["createExperimentByTemplate"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/templates/{id}/experiment-execute": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Create an experiment based on an experiment template and run experiment
         * @description Use the given experiment template id and the placeholder values to create or update and immediately run the experiment. The `externalId` will be used to identify whether the experiment exists already and should be updated or newly inserted.
         */
        post: operations["saveAndRunFromTemplate"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/templates/{id}/experiment-update/{key}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Update an existing experiment based on a template
         * @description Use the given experiment template id and the placeholder values to create or update the experiment. Placeholders that have been used for the initial creation will be reused. Provided placeholders from the body will overwrite existing placeholders.
         */
        post: operations["updateExperimentByTemplate"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/experiments/templates/imports": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Import experiment templates
         * @description Import experiment templates with given IDs from linked hub.
         */
        post: operations["importFromHub"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/explore/landscape/views": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch all saved landscape views of a team
         * @description Get a list of all saved explorer landscape views that belong to the given team.
         */
        get: operations["getLandscapeViews"];
        put?: never;
        /**
         * Create a saved landscape view
         * @description Create a new saved explorer landscape view for a team.
         */
        post: operations["createLandscapeView"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/explore/landscape/views/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single saved landscape view
         * @description Get all details of a single saved explorer landscape view.
         */
        get: operations["getLandscapeView"];
        /**
         * Update a saved landscape view
         * @description Update an existing saved explorer landscape view. The view's thumbnail is rendered by the UI and cannot be produced through the API. When an update changes a field that affects how the landscape renders (environment, filter query, group-by, size-by, color-by or show-advice), the thumbnail is cleared so it is not left stale; a metadata-only change (e.g. name or description) keeps it. The UI regenerates the thumbnail on its next save.
         */
        put: operations["updateLandscapeView"];
        post?: never;
        /**
         * Delete a saved landscape view
         * @description Remove the given saved explorer landscape view from the Steadybit platform.
         */
        delete: operations["deleteLandscapeView"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/health": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get: operations["health"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/health/liveness": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get: operations["liveness"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/health/readiness": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get: operations["readiness"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/hubs": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of all hubs
         * @description Get a list of all hubs that are currently connected.
         */
        get: operations["getHubs"];
        put?: never;
        /**
         * Create or update a hub
         * @description Insert or update a hub. The `id` will be used to identify whether the hub exists already and should be updated or newly inserted. The hub content can be synchronized depending on the `resync` parameter. This operation is executed synchronously and may take some time to complete depending on the network connection to hub address.
         */
        post: operations["upsertHub"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/hubs/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single hub
         * @description Get all details of a single hub.
         */
        get: operations["getHubById"];
        put?: never;
        post?: never;
        /**
         * Delete a hub
         * @description Remove the given hub.
         */
        delete: operations["deleteHub"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/hubs/{id}/resync": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Re-synchronize a hub
         * @description Fetch the latest hub definition based on `hubRepository`. This operation is executed synchronously and may take some time to complete depending on the network connection to hub address.
         */
        post: operations["resyncHub"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/hubs/connection-check": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Check a hub connection
         * @description Check if the given hub connection details point to a valid hub.
         */
        post: operations["connectionCheck"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/integrations/preflight": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of preflight webhooks
         * @description Get a list of all existing preflight webhooks.
         */
        get: operations["getPreflightWebhooks"];
        put?: never;
        /**
         * Create or update a preflight webhook
         * @description Insert or update a preflight webhook.<br/This endpoint requires an admin-token and can't be used with a team-based token.
         */
        post: operations["upsertPreflightWebhook"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/integrations/preflight-action": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of preflight action integrations
         * @description Get a list of all existing preflight preflight action integrations.
         */
        get: operations["getPreflightActionIntegrations"];
        put?: never;
        /**
         * Create or update a preflight action integration
         * @description Insert or update a preflight action integration.<br/This endpoint requires an admin-token and can't be used with a team-based token.
         */
        post: operations["upsertPreflightActionIntegration"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/integrations/preflight-action/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single preflight preflight action integrations
         * @description Get all details of a single preflight preflight action integrations.
         */
        get: operations["getPreflightActionIntegration"];
        put?: never;
        post?: never;
        /**
         * Delete preflight action integration
         * @description Remove the given preflight action integration from the Steadybit platform.
         */
        delete: operations["deletePreflightActionIntegration"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/integrations/preflight/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single preflight webhooks
         * @description Get all details of a single preflight webhooks.
         */
        get: operations["getPreflightWebhook"];
        put?: never;
        post?: never;
        /**
         * Delete preflight webhook
         * @description Remove the given preflight webhook from the Steadybit platform.
         */
        delete: operations["deletePreflightWebhook"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/integrations/slack": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of slack integrations
         * @description Get a list of all existing slack integrations.
         */
        get: operations["getSlackIntegrations"];
        put?: never;
        /**
         * Create or update a slack integration
         * @description Insert or update a slack integration.<br/This endpoint requires an admin-token and can't be used with a team-based token.
         */
        post: operations["upsertSlackIntegration"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/integrations/slack/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single slack integration
         * @description Get all details of a single slack integration.
         */
        get: operations["getSlackIntegration"];
        put?: never;
        post?: never;
        /**
         * Delete slack integration
         * @description Remove the given slack integration from the Steadybit platform.
         */
        delete: operations["deleteSlackIntegration"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/integrations/webhook": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of custom webhooks
         * @description Get a list of all existing custom webhooks.
         */
        get: operations["getCustomWebhooks"];
        put?: never;
        /**
         * Create or update a custom webhook
         * @description Insert or update a custom webhook.<br/This endpoint requires an admin-token and can't be used with a team-based token.
         */
        post: operations["upsertCustomWebhook"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/integrations/webhook/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single custom webhooks
         * @description Get all details of a single custom webhooks.
         */
        get: operations["getCustomWebhook"];
        put?: never;
        post?: never;
        /**
         * Delete custom webhook
         * @description Remove the given custom webhook from the Steadybit platform.
         */
        delete: operations["deleteCustomWebhook"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/killswitch": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get the current status of the kill switch
         * @description Determines the current status of the kill switch without changing it.
         */
        get: operations["getKillswitch"];
        put?: never;
        /**
         * Activate / engage the kill switch
         * @description Activates / engages the kill switch to cancel all experiments running at the moment and prevent execution of new experiments until the kill switch is disengaged / deactivated again.
         */
        post: operations["engageKillswitch"];
        /**
         * Deactivate / disengage the kill switch
         * @description Deactivates / disengages the kill switch to allow experiment runs.<br/>Experiment runs that were not executed due to engaged / active kill switch will not be automatically executed, they need to be triggered again.
         */
        delete: operations["disengageKillswitch"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/license": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get license summary. */
        get: operations["getLicenseSummary"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/license/report": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get license report. */
        get: operations["getReport"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/preflight/actions": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get all preflight actions. */
        get: operations["getPreflightActionSummary"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/properties/associations": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get all current associations. */
        get: operations["getAssociations"];
        put?: never;
        /**
         * Create or update a property association
         * @description Insert or update the property association. The `id` will be used to identify whether the schedule exists already and should be updated or newly inserted.
         *
         *        Examples:
         *        - Assign the property `RESULT_COLOR` to all experiment designs:
         *        ```
         *        {
         *            "key": "RESULT_COLOR",
         *            "editableInExecution": false,
         *            "required": true
         *        }
         *        ```
         *        - Assign the property `RESULT_COLOR` to the design ADM-15:
         *        ```
         *        {
         *            "key": "RESULT_COLOR",
         *            "editableInExecution": false,
         *            "experimentKey": "ADM-15",
         *            "required": true
         *        }
         *        ```
         *        - Assign the property `RESULT_COLOR` that can be edited in each experiment execution of the experiment with key `ADM-15`:
         *        ```
         *        {
         *            "key": "RESULT_COLOR",
         *            "editableInExecution": true,
         *            "experimentKey": "ADM-15",
         *            "required": false
         *        }
         *        ```
         *        - Assign the property `RESULT_COLOR` to a service `0a2d67b9-1d5a-4179-8c32-e5296be1f56f`:
         *        ```
         *        {
         *            "key": "RESULT_COLOR",
         *            "associationType": "SERVICE",
         *            "serviceId": "0a2d67b9-1d5a-4179-8c32-e5296be1f56f",
         *            "required": false
         *        }
         *        ```
         */
        post: operations["upsertPropertyAssociation"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/properties/associations/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get property association by a given id. */
        get: operations["getPropertyDefinition_1"];
        put?: never;
        post?: never;
        /** Remove an existing property association. */
        delete: operations["deletePropertyAssociation"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/properties/definitions": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get: operations["getPropertyDefinitions"];
        put?: never;
        /**
         * Create or update property definition
         * @description Insert or update the property definition specified by the given `key`.
         */
        post: operations["upsertPropertyDefinition"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/properties/definitions/{key}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get property definition for a specific property definition key. */
        get: operations["getPropertyDefinition"];
        put?: never;
        post?: never;
        /** Remove an existing property definition */
        delete: operations["deletePropertyDefinition"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/reports/environments": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Get environment counts over time
         * @description Returns the number of environments in the tenant aggregated into time buckets.
         */
        post: operations["getEnvironmentCounts"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/reports/experiments/created": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Get experiment creation counts over time
         * @description Returns experiment creation counts aggregated into time buckets, optionally grouped by creation method or origin.
         */
        post: operations["getExperimentCreations"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/reports/experiments/executed": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Get experiment execution counts over time
         * @description Returns experiment execution counts aggregated into time buckets, optionally grouped by state, trigger, or attack action.
         */
        post: operations["getExperimentExecutions"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/reports/services/average": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Get average service risk over time
         * @description Returns the average risk across services aggregated into time buckets. Risk is reported as an integer 0-100.
         */
        post: operations["getAverageRisk"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/reports/services/by-category": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Get average service risk grouped by category over time
         * @description Returns the average risk per category (key from the categoryRisks map) across services, aggregated into time buckets.
         */
        post: operations["getRiskByCategory"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/reports/services/distribution": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Get service risk level distribution over time
         * @description Returns the count of services in each risk level (LOW, MEDIUM, HIGH) aggregated into time buckets.
         */
        post: operations["getRiskDistribution"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/reports/teams": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Get team counts over time
         * @description Returns the number of teams in the tenant aggregated into time buckets.
         */
        post: operations["getTeamCounts"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/reports/users": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Get user counts over time
         * @description Returns the number of users in the tenant aggregated into time buckets.
         */
        post: operations["getUserCounts"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/services": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Fetch a list of services */
        get: operations["getServiceList"];
        put?: never;
        /**
         * Create or update service
         * @description Insert or update the service specified by the given `id`.
         */
        post: operations["upsertService"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/services/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get service by id. */
        get: operations["getService"];
        put?: never;
        post?: never;
        /** Delete an existing service */
        delete: operations["deleteService"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/services/{id}/experiments": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get experiments associated to an service. */
        get: operations["getServiceExperiments"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/services/{id}/experiments/custom": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Link a custom experiment to a service. */
        post: operations["linkCustomExperiment"];
        /** Remove a linked custom experiment from a service. */
        delete: operations["unlinkCustomExperiment"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/services/{id}/experiments/provided": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /** Create or update a provided experiment. */
        post: operations["upsertProvidedExperiment"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/services/{id}/risk": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get the risk score for a service */
        get: operations["getRisk"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/services/{id}/variables": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get service variables
         * @description Get all variables owned by the service.
         */
        get: operations["getServiceVariables"];
        /**
         * Replace all service variables
         * @description All provided variables will be associated with the given service and existing ones removed.
         */
        put: operations["setServiceVariables"];
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Add / merge service variables
         * @description All provided variables will be associated with the given service.<br/>If a variable key is already in use, its value is updated.<br/>If a variable is already associated but not provided, it continues to exist.
         */
        patch: operations["mergeServiceVariables"];
        trace?: never;
    };
    "/api/services/profiles": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Fetch a list of service profiles */
        get: operations["getProfiles"];
        put?: never;
        /**
         * Create or update service profile
         * @description Insert or update the service profile specified by the given `id`.
         */
        post: operations["upsertProfile"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/services/profiles/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Get service profile by id. */
        get: operations["getProfile"];
        put?: never;
        post?: never;
        /** Delete an existing service profile */
        delete: operations["deleteProfile"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/target-stats": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /** Gather target statistics without any filters */
        get: operations["getTargetsStats"];
        put?: never;
        /** Gather target statistics for a given predicate or query */
        post: operations["getTargetsStats_1"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/targets": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get targets
         * @description Get targets.
         */
        get: operations["getTargets"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/targets/attributes/keys": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get attribute key
         * @description Get all available attribute keys for a specific target type in a given environment.
         */
        get: operations["getTargetAttributeKeys"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/targets/attributes/values": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get attribute values
         * @description Get all available attribute values for a specific attribute and target type in a given environment.
         */
        get: operations["getTargetAttributeValues"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/teams": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a list of all teams
         * @description Get a list of all teams that exist.<br/> If used with a team-associated `accessToken` and `onlyAccessible` is set to `true` you only get the team of the `accessToken`.
         */
        get: operations["getTeams"];
        put?: never;
        /**
         * Create or update a team
         * @description Insert or update the team in Steadybit. The `key` will be used to identify whether the team exists already and should be updated or newly inserted. If a provided member's username or email is not yet known it will be skipped.
         */
        post: operations["upsertTeam"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/teams/{key}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Fetch a single team
         * @description Get all details of a single existing teams.
         */
        get: operations["getTeam"];
        put?: never;
        post?: never;
        /**
         * Delete team
         * @description Remove the given team from the Steadybit platform. This will only work, if there are no experiments running at the moment.
         */
        delete: operations["deleteTeam"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/teams/{key}/environments": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get all environments assigned to the team
         * @description Get a list of members that are part of the specified team. The list contains the username, being a Steadybit user id, and the role in this particular team (owner or member).
         */
        get: operations["getTeamEnvironments"];
        /**
         * Update the environments of a specific team
         * @description The allowed environments of the specified team will be updated with these provided. New environments will be added to the team, environments not provided in the request will be removed from the team.
         */
        put: operations["setTeamEnvironments"];
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/teams/{key}/environments/add": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Add an allowed environment to a team
         * @description The given environments will be added to the specified team.
         */
        post: operations["addTeamEnvironments"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/teams/{key}/environments/remove": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Remove allowed environment from a team
         * @description The given environments will be removed from the specified team.
         */
        post: operations["removeTeamEnvironments"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/teams/{key}/members": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Get all members being part of the team
         * @description Get a list of members that are part of the specified team. The list contains the username, being a Steadybit user id, and the role in this particular team (owner or member).
         */
        get: operations["getTeamMembers"];
        /**
         * Update the members of a specific team
         * @description The members of the specified team will be updated with these provided. New team members will be added to the team, team members not provided in the request will be removed from the team.
         */
        put: operations["setTeamMembers"];
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/teams/{key}/members/add": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Add team members to a team
         * @description The given members will be added to the specified team.
         */
        post: operations["addTeamMembers"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/teams/{key}/members/remove": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Remove team members from a team
         * @description The given members will be removed from the specified team. However, they are still able to login, view the content of the team and may still be member of another team.
         */
        post: operations["removeTeamMembers"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/api/users/invite": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Invite users to a tenant
         * @description Invite users to a tenant. The invited users will receive an email with an invitation link to join the tenant.
         */
        post: operations["inviteUser"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        /**
         * @description A step that is executed as part of an experiment.
         * @example [
         *       {
         *         "id": "40b0f797-912d-4256-8887-1553561962a9",
         *         "predecessorId": null,
         *         "ignoreFailure": false,
         *         "parameters": {
         *           "duration": "10s"
         *         }
         *       }
         *     ]
         */
        AbstractExperimentExecutionStepAO: {
            /**
             * @description Custom label assigned during experiment design to express the intention of this step
             * @example Container 'xyz' can not be reached
             */
            customLabel?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * Format: uuid
             * @description Unique identifier of this step execution
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            id?: string;
            /**
             * @description Whether the experiment should fail/error immediately in case this step fails/errors.
             * @example false
             */
            ignoreFailure?: boolean;
            /**
             * @description Step-specific parameters of the experiment step configuration
             * @example {
             *       "duration": "10s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            /**
             * Format: uuid
             * @description Unique identifier of the step execution that precedes this step, null if it is the first step of a lane
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            predecessorId?: string;
            /**
             * @description Reason in case this experiment step execution failed or errored
             * @example Couldn't read state of container...
             */
            reason?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step was started
             * @example 2023-01-01T09:00:00Z
             */
            started?: string;
            /**
             * @description Current state of this step in the experiment (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
            /**
             * @description Type of this step execution (e.g. ACTION, WAIT)
             * @example ACTION
             */
            stepType: string;
        } & (components["schemas"]["ExperimentExecutionStepActionAO"] | components["schemas"]["ExperimentExecutionStepWaitAO"] | components["schemas"]["ExperimentExecutionStepServiceValidationAO"]);
        /**
         * @description A step that is executed as part of an experiment.
         * @example [
         *       {
         *         "id": "40b0f797-912d-4256-8887-1553561962a9",
         *         "predecessorId": null,
         *         "ignoreFailure": false,
         *         "parameters": {
         *           "duration": "10s"
         *         }
         *       }
         *     ]
         */
        AbstractWebhookPayloadExecutionStepAO: {
            /**
             * @description Custom label assigned during experiment design to express the intention of this step
             * @example Container 'xyz' can not be reached
             */
            customLabel?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * Format: uuid
             * @description Unique identifier of this step execution
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            id?: string;
            /**
             * @description Whether the experiment should fail/error immediately in case this step fails/errors.
             * @example false
             */
            ignoreFailure?: boolean;
            /**
             * @description Step-specific parameters of the experiment step configuration
             * @example {
             *       "duration": "10s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            /**
             * Format: uuid
             * @description Unique identifier of the step execution that precedes this step, null if it is the first step of a lane
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            predecessorId?: string;
            /**
             * @description Reason in case this experiment step execution failed or errored
             * @example Couldn't read state of container...
             */
            reason?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step was started
             * @example 2023-01-01T09:00:00Z
             */
            started?: string;
            /**
             * @description Current state of this step in the experiment (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
        } & (components["schemas"]["WebhookPayloadExecutionStepWaitAO"] | components["schemas"]["WebhookPayloadExecutionStepActionAO"] | components["schemas"]["WebhookPayloadExecutionStepServiceValidationAO"]);
        /**
         * @description The logged event was performed via API authorized via access token
         * @example {
         *       "id": "VDKTEBLl",
         *       "name": "CI/CD",
         *       "tokenType": "TEAM",
         *       "principalType": "ACCESS_TOKEN"
         *     }
         */
        AccessTokenPrincipalAL: {
            /**
             * @description Unique identifier of this access token principal
             * @example VDKTEBLl
             */
            id: string;
            /**
             * @description Name of the access token that was used
             * @example CI/CD
             */
            name: string;
            /**
             * @description Principal type for access token based principal
             * @example ACCESS_TOKEN
             * @enum {string}
             */
            principalType: "USER" | "ACCESS_TOKEN" | "BATCH_JOB";
            /**
             * @description Access token type that was used to perform the logged event
             * @example TEAM
             * @enum {string}
             */
            tokenType: "ADMIN" | "TEAM" | "WILDCARD";
        };
        /**
         * @description A single access token of a team. The token itself can't be read again
         * @example {
         *       "id": "aP4cDVfA",
         *       "name": "CI/CD"
         *     }
         */
        AccessTokensPageItemAO: {
            /**
             * @description Unique identifier of the access token
             * @example CQer2Oar
             */
            id?: string;
            /**
             * @description Name of the Access Token to document e.g. its purpose
             * @example CI/CD access token
             */
            name?: string;
            /**
             * @description Team associated with this token or null if this is an admin token
             * @example ADM
             */
            team?: string | null;
            /**
             * @description Type of the access token
             * @example ADMIN
             * @enum {string}
             */
            type?: "ADMIN" | "TEAM";
        };
        /** @description A single access token */
        AccessTokensPageItemV2AO: {
            /**
             * Format: date-time
             * @description Expiration date of the token. Null means the token never expires.
             */
            expiresAt?: string;
            /** @description Unique identifier of this access token */
            id?: string;
            /**
             * Format: date-time
             * @description Date of the last token usage. Null means the token was never used.
             */
            lastUsed?: string;
            /** @description Name of this access token */
            name?: string;
            /** @description Teams associated with this token. */
            teams?: string[];
            /**
             * @description Type of this token
             * @enum {string}
             */
            type?: "ADMIN" | "TEAM" | "WILDCARD";
        };
        /** @description An action that is currently registered. */
        ActionAO: {
            /**
             * @description Category grouping similar actions.
             * @example Resource
             */
            category?: string;
            defaultBlastRadius: components["schemas"]["DefaultBlastRadiusAO"];
            /** @description Description of the action. */
            description: string;
            hint?: components["schemas"]["HintAO"];
            hubSummary?: string | null;
            /** @description Icon of the action as a data URI (may be a large base64-encoded image). */
            icon?: string | null;
            /**
             * @description Unique identifier of the action.
             * @example com.steadybit.extension_container.stress_cpu
             */
            id: string;
            /**
             * @description Kind of the action.
             * @example ATTACK
             * @enum {string}
             */
            kind: "ATTACK" | "CHECK" | "LOAD_TEST" | "OTHER" | "BASIC";
            /** @description Parameters that describe how to fetch metrics for this action. */
            metricQueryParameters: components["schemas"]["ParameterAO"][];
            /**
             * @description Behavior when the target query does not match any targets.
             * @example INCLUDE_NONE
             * @enum {string}
             */
            missingQuerySelection: "INCLUDE_ALL" | "INCLUDE_NONE";
            /**
             * @description Display name of the action.
             * @example Stress CPU
             */
            name: string;
            parameters?: components["schemas"]["ParameterAO"][];
            /**
             * @description Restriction on the number of targets this action may operate on.
             * @example NONE
             * @enum {string}
             */
            quantityRestriction: "EXACTLY_ONE" | "ALL" | "NONE";
            /** @description Whether this action supports metric queries. */
            supportsMetricQueries?: boolean;
            target?: components["schemas"]["TargetSelectorAO"];
            /** @description Predefined target predicate templates offered to the user when configuring this action. */
            targetPredicateTemplates: components["schemas"]["TargetPredicateTemplateAO"][];
            /**
             * @description Technology the action belongs to.
             * @example Container
             */
            technology?: string;
            /**
             * @description Version of the action.
             * @example 1.2.3
             */
            version?: string;
        };
        /** @description List of actions. */
        ActionSummariesAO: {
            /** @description List of actions. */
            actions?: components["schemas"]["ActionAO"][];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        /**
         * @description Add a value to a list property of an execution
         * @example {
         *       "type": "add_value_to_list_property",
         *       "propertyKey": "observations",
         *       "value": "This looks interesting!"
         *     }
         */
        AddValueToListProperty: {
            type: "AddValueToListProperty";
        } & (Omit<components["schemas"]["ExecutionModification"], "type"> & {
            /**
             * @description The key of the property.
             * @example observations
             */
            propertyKey: string;
            /**
             * @description The value that should be added to the list. Number in case of a numeric list, String otherwise
             * @example This looks interesting!
             */
            value: Record<string, never>;
        });
        AdvancedRadiusAO: {
            /**
             * @description The target attribute that should be picked randomly
             * @example aws.zone
             */
            attribute: string;
            /**
             * @description Only in execution - the values that has been picked by the randomizer for the given execution
             * @example ['us-east-1a','us-east-1b']
             */
            pickedValues?: string[] | null;
            /**
             * @description The percentage (example: `50%`) or fixed amount (example: `15#`)
             * @example 50%
             */
            value: string;
        };
        /**
         * @description A pageable list of pieces of advice.
         * @example {
         *       "totalItems": 108,
         *       "nextOffset": 3,
         *       "items": [
         *         {
         *           "target": {
         *             "type": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *             "reference": "prod-demo/steadybit-demo/gateway",
         *             "label": "gateway"
         *           },
         *           "advice": {
         *             "type": "com.steadybit.extension_kubernetes.advice.k8s-cpu-limit",
         *             "label": "Limit CPU Resources",
         *             "tags": [
         *               "kubernetes",
         *               "limit",
         *               "cpu"
         *             ],
         *             "status": "Validation needed",
         *             "summary": "You already took action and configured a CPU limit. Validate your configuration via an experiment."
         *           },
         *           "url": "https://platform.steadybit.com/permalink/advice/eyAiZW52..."
         *         },
         *         {
         *           "target": {
         *             "type": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *             "reference": "prod-demo/steadybit-demo/gateway",
         *             "label": "gateway"
         *           },
         *           "advice": {
         *             "type": "com.steadybit.extension_kubernetes.advice.k8s-cpu-request",
         *             "label": "Requesting Reasonable CPU Resources",
         *             "tags": [
         *               "kubernetes",
         *               "request",
         *               "cpu"
         *             ],
         *             "status": "Validation needed",
         *             "summary": "You specified a CPU request that informs Kubernetes decision where to schedule your pods of *activemq*.\nPlease confirm that your requested CPU share is reasonable for your type of application."
         *           },
         *           "url": "https://platform.steadybit.com/permalink/advice/eyAiZW52..."
         *         }
         *       ]
         *     }
         */
        AdviceSummaryAO: {
            items?: components["schemas"]["TargetAdviceAO"][];
            /**
             * Format: int32
             * @description Next queryable offset to query for next batch of advice
             * @example 21
             */
            nextOffset?: number | null;
            /**
             * Format: int64
             * @description Total amount of advice matching your query
             * @example 241
             */
            totalItems?: number;
        };
        /**
         * @description An attributes (key-value-pair) that is associated to a target
         * @example {
         *       "key": "container.port",
         *       "value": "51152:2376"
         *     }
         */
        Attribute: {
            /**
             * @description The key of the attribute, may be associated multiple times to the same target
             * @example container.engine
             */
            key: string;
            /**
             * @description The value of the attribute
             * @example docker
             */
            value: string;
        };
        /**
         * @description An attributes (key-value-pair) that is associated to a target
         * @example {
         *       "key": "container.port",
         *       "value": "51152:2376"
         *     }
         */
        AttributeAO: {
            /**
             * @description The key of the attribute, may be associated multiple times to the same target
             * @example container.engine
             */
            key: string;
            /**
             * @description The value of the attribute
             * @example docker
             */
            value: string;
        };
        /**
         * @description Audit log entry.
         * @example {
         *       "id": "14av1421-aol3-4159-8ae2-47f5a9ba119e",
         *       "tenant": {
         *         "key": "Demo",
         *         "name": "Demo Tenant"
         *       },
         *       "trigger": {
         *         "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
         *         "triggerType": "HTTP_REQUEST"
         *       },
         *       "eventName": "experiment.created",
         *       "eventTime": "2023-01-03T09:13:00Z",
         *       "experiment": {
         *         "key": "ADM-1"
         *       }
         *     }
         */
        AuditLogEntry: {
            environment?: components["schemas"]["EnvironmentAL"];
            /**
             * @description Event name that was audited
             * @example experiment.created
             */
            eventName: string;
            /**
             * Format: date-time
             * @description The time at which the event was audited
             * @example 2023-01-03T09:13:00Z
             */
            eventTime: string;
            /**
             * Format: uuid
             * @description Unique identifier of the audit log entry
             */
            id: string;
            principal?: components["schemas"]["PrincipalAL"];
            team?: components["schemas"]["TeamAL"];
            tenant: components["schemas"]["TenantAL"];
            trigger?: components["schemas"]["AuditLogTrigger"];
        };
        /**
         * @description The trigger that caused the event to happen
         * @example {
         *       "userAgent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/112.0.0.0 Safari/537.36",
         *       "triggerType": "HTTP_REQUEST"
         *     }
         */
        AuditLogTrigger: {
            triggerType: string;
        } | null;
        /**
         * @description A single step in a lane.
         * @example {
         *       "type": "action",
         *       "ignoreFailure": false,
         *       "parameters": {
         *         "duration": "30s"
         *       },
         *       "actionType": "com.steadybit.extension_host.stress-cpu",
         *       "radius": {
         *         "targetType": "com.steadybit.extension_container.container",
         *         "predicate": {
         *           "operator": "AND",
         *           "predicates": [
         *             {
         *               "key": "k8s.deployment",
         *               "operator": "EQUALS",
         *               "values": [
         *                 "hot-deals"
         *               ]
         *             }
         *           ]
         *         },
         *         "query": null,
         *         "percentage": 100
         *       }
         *     }
         */
        BaseExperimentStepAO: {
            customLabel?: string;
            /**
             * @description Ignore any errors and failures of this single step and continue the execution of an experiment run
             * @example false
             */
            ignoreFailure?: boolean;
            /** @description Optional metric checks used to define success or failure of this step */
            metricChecks?: components["schemas"]["MetricCheckAO"][];
            /** @description Optional metric queries used of this step to filter e.g. monitoring data */
            metricQueries?: components["schemas"]["MetricQueryAO"][];
            /**
             * @description Configuration parameters of this step that are saved during experiment design and evaluated at execution time.
             * @example {
             *       "duration": "30s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            type: string;
        } & (components["schemas"]["ExperimentStepActionAO"] | components["schemas"]["ExperimentStepWaitAO"] | components["schemas"]["ExperimentStepServiceValidationAO"]);
        /**
         * @description A batch job has performed the logged event
         * @example {
         *       "username": "af1bw7kj-d299-47ab-998f-c2a53b433820",
         *       "principalType": "BATCH_JOB"
         *     }
         */
        BatchPrincipalAL: {
            /**
             * @description Principal type for batch based principal
             * @example BATCH_JOB
             * @enum {string}
             */
            principalType: "USER" | "ACCESS_TOKEN" | "BATCH_JOB";
            /**
             * @description Username of the user, internal identifier of Steadybit
             * @example 13av2737-b318-4048-a79d-4789d645bc31
             */
            username?: string;
        };
        /**
         * @description Blast radius that is applied to define the set of targets as well as an optional random subset
         * @example {
         *       "targetType": "com.steadybit.extension_container.container",
         *       "percentage": 50,
         *       "predicate": {
         *         "operator": "AND",
         *         "predicates": [
         *           {
         *             "key": "container.host/name",
         *             "operator": "EQUALS",
         *             "values": [
         *               "docker-desktop/minikube"
         *             ]
         *           }
         *         ]
         *       }
         *     }
         */
        BlastRadiusAO: {
            /**
             * Format: int32
             * @description In case a fixed number of as subset of specified targets should be effected
             * @example 2
             */
            maximum?: number;
            /**
             * Format: int32
             * @description In case a percentage subset of the specified targets should be effected
             * @example 40
             */
            percentage?: number;
            predicate?: components["schemas"]["TargetPredicateAO"];
            /**
             * @description Target type that is effected by that action
             * @example container
             */
            targetType?: string;
        };
        /**
         * @description The risk for a given category in a service
         * @example {
         *       "total": 50,
         *       "experiment": 50,
         *       "advice": 50
         *     }
         */
        CategoryRiskAO: {
            /**
             * Format: int32
             * @description The advice risk for this category, or null when no advice could be found for the given target selection
             */
            advice?: number;
            /**
             * Format: int32
             * @description The experiment risk for this category
             */
            experiment?: number;
            /**
             * Format: int32
             * @description The total risk for this category
             */
            total?: number;
        };
        ComparableValueAO: {
            type?: string;
        };
        CreateAccessTokenRequestAO: {
            /**
             * @description Name of the Access Token to document its purpose
             * @example CI/CD access token
             */
            name: string;
            /**
             * @description Team associated with this token or null if this is an admin token
             * @example ADM
             */
            team?: string | null;
            /**
             * @description Type of this token.
             * @example TEAM
             * @enum {string}
             */
            type: "ADMIN" | "TEAM";
        };
        CreateAccessTokenRequestV2AO: {
            /**
             * Format: date-time
             * @description Expiration date of the token. If not set, the token will never expire.
             * @example 2027-01-01T00:00:00Z
             */
            expiresAt?: string | null;
            /**
             * @description Name of the Access Token to document its purpose
             * @example CI/CD access token
             */
            name: string;
            /**
             * @description Keys of teams to associate with this token. Required when type TEAM, must be empty for type ADMIN.
             * @example [
             *       "ADM",
             *       "DEV"
             *     ]
             */
            teams?: string[] | null;
            /**
             * @description Type of this token.
             * @example TEAM
             * @enum {string}
             */
            type: "ADMIN" | "TEAM" | "WILDCARD";
        };
        CreateAccessTokenResponseAO: {
            /**
             * @description Unique identifier of the access token
             * @example CQer2Oar
             */
            id?: string;
            /**
             * @description Token to be used to authenticate in the API. <br/> Make sure to save the generated token as you can't read it again afterwards for security-reasons.
             * @example a1fDXcA0.P.2dFfGl3fAq126mnVCxyPZLoEmLwPi2
             */
            token?: string;
        };
        CreateAccessTokenResponseV2AO: {
            /** @description Unique identifier of this access token */
            id?: string;
            /** @description The access token. Make sure to save it as you can't read it again afterwards for security-reasons. */
            token?: string;
        };
        /**
         * @description Create or update the experiment with the given experiment design.
         * @example {
         *       "name": "Blackhole Hot-deals",
         *       "team": "ADM",
         *       "environment": "Global",
         *       "lanes": [
         *         {
         *           "steps": [
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "30s"
         *               },
         *               "actionType": "com.steadybit.extension_host.stress-cpu",
         *               "radius": {
         *                 "targetType": "com.steadybit.extension_container.container",
         *                 "predicate": {
         *                   "operator": "AND",
         *                   "predicates": [
         *                     {
         *                       "key": "k8s.deployment",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "hot-deals"
         *                       ]
         *                     }
         *                   ]
         *                 },
         *                 "query": null,
         *                 "percentage": 100
         *               }
         *             }
         *           ]
         *         }
         *       ],
         *       "properties": {
         *         "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
         *       }
         *     }
         */
        CreateAndRunExperimentAO: {
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment?: string;
            /**
             * @description Variables that will be merged for the single experiment execution with the variables defined of the experiment or environment that the experiment will be executed in. A `key` that exists already in the experiment or environment variables will be overridden for this execution, all others will be added solely in the context of the first experiment execution. Each value is either a constant string, an array of constant strings, or a select expression object.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal"
             *     }
             */
            executionVariables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * @description Variables that will be used when the experiment will be executed. Experiment variables will override existing environment variables. Each value is either a constant string, an array of constant strings, or a select expression object (`{"type":"select",...}`).
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal",
             *       "targetServices": [
             *         "gateway",
             *         "hot-deals",
             *         "fashion-bestseller"
             *       ],
             *       "httpEndpointZones": {
             *         "type": "select",
             *         "targetType": "com.steadybit.extension_container.container",
             *         "attribute": "aws.zone",
             *         "mode": "fixed",
             *         "count": 2
             *       }
             *     }
             */
            experimentVariables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * @description An optional external identifier used for create-or-update semantics.
             * @example 1234567
             */
            externalId?: string;
            /**
             * @deprecated
             * @description An optional external reference. Will be removed and is replaced by tags. If used with experiment creation, the value will be added as a tag.
             * @example INCIDENT-4711
             */
            externalReference?: string;
            /**
             * @description The hypothesis that is validated by the experiment
             * @example System is able to survive a latency in the network of 1500ms
             */
            hypothesis?: string;
            /**
             * @description The lanes (steps executed in parallel) in the experiment. Each lane consists of multiple steps that are executed sequential per lane.
             * @example [
             *       {
             *         "steps": [
             *           {
             *             "type": "action",
             *             "ignoreFailure": false,
             *             "parameters": {
             *               "duration": "30s"
             *             },
             *             "actionType": "com.steadybit.extension_host.stress-cpu",
             *             "radius": {
             *               "targetType": "com.steadybit.extension_container.container",
             *               "predicate": {
             *                 "operator": "AND",
             *                 "predicates": [
             *                   {
             *                     "key": "k8s.deployment",
             *                     "operator": "EQUALS",
             *                     "values": [
             *                       "hot-deals"
             *                     ]
             *                   }
             *                 ]
             *               },
             *               "query": null,
             *               "percentage": 100
             *             }
             *           }
             *         ]
             *       }
             *     ]
             */
            lanes: components["schemas"]["ExperimentLaneAO"][];
            /**
             * @description Name of the experiment to easily identify the experiment
             * @example Shop survives unavailability of hot-deals products
             */
            name: string;
            /**
             * @description The properties of the experiment
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * @description Team keys with which the experiment is shared with
             * @example [OPS, SHOP]
             */
            sharedTeams?: string[];
            /**
             * @description An optional set of tags you can use to search for.
             * @example [
             *       "myTag",
             *       "myOtherTag"
             *     ]
             */
            tags?: string[];
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            team: string;
        };
        /**
         * @description Create or update an experiment based on an experiment template and run it.
         * @example {
         *       "environment": "steadybit-demo",
         *       "team": "DEMO",
         *       "placeholders": [
         *         {
         *           "key": "CLUSTER",
         *           "value": "demo-cluster"
         *         },
         *         {
         *           "key": "BOOL",
         *           "value": true
         *         },
         *         {
         *           "key": "NUMBER",
         *           "value": 15
         *         },
         *         {
         *           "key": "KEYVALUE",
         *           "value": [
         *             {
         *               "key": "example-a",
         *               "value": "abc"
         *             },
         *             {
         *               "key": "example-b",
         *               "value": "123"
         *             }
         *           ]
         *         },
         *         {
         *           "key": "LIST",
         *           "value": [
         *             "entry1",
         *             "entry2",
         *             "entry3"
         *           ]
         *         },
         *         {
         *           "key": "FILE",
         *           "value": {
         *             "fileName": "example.txt",
         *             "data": "SGVsbG8gV29ybGQh"
         *           }
         *         }
         *       ],
         *       "externalId": "1234567",
         *       "executionVariables": {
         *         "httpEndpoint": "http://dev.shop.products.internal"
         *       }
         *     }
         */
        CreateAndRunExperimentFromTemplateAO: {
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment?: string;
            /**
             * @description Variables that will be merged for the single experiment execution with the variables defined of the experiment or environment that the experiment will be executed in. A `key` that exists already in the experiment or environment variables will be overridden for this execution, all others will be added solely in the context of the first experiment execution. Each value is either a constant string, an array of constant strings, or a select expression object.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal"
             *     }
             */
            executionVariables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * @description Variables that will be added to the created experiment design. A `key` that exists already in the environment variables will be overridden. Each value is either a constant string, an array of constant strings, or a select expression object.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal"
             *     }
             */
            experimentVariables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * @description An optional external identifier used for create-or-update semantics.
             * @example 1234567
             */
            externalId?: string;
            /** @description List of template placeholder values */
            placeholders?: components["schemas"]["ExperimentTemplatePlaceholderValueAO"][];
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            team: string;
        };
        /**
         * @description Create or update the experiment with the given experiment design.
         * @example {
         *       "name": "Blackhole Hot-deals",
         *       "team": "ADM",
         *       "environment": "Global",
         *       "sharedTeams": [
         *         "SHOP"
         *       ],
         *       "lanes": [
         *         {
         *           "steps": [
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "30s"
         *               },
         *               "actionType": "com.steadybit.extension_host.stress-cpu",
         *               "radius": {
         *                 "targetType": "com.steadybit.extension_container.container",
         *                 "predicate": {
         *                   "operator": "AND",
         *                   "predicates": [
         *                     {
         *                       "key": "k8s.deployment",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "hot-deals"
         *                       ]
         *                     }
         *                   ]
         *                 },
         *                 "query": null,
         *                 "percentage": 100
         *               }
         *             }
         *           ]
         *         }
         *       ],
         *       "properties": {
         *         "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
         *       }
         *     }
         */
        CreateExperimentAO: {
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment?: string;
            /**
             * @description Variables that will be used when the experiment will be executed. Experiment variables will override existing environment variables. Each value is either a constant string, an array of constant strings, or a select expression object (`{"type":"select",...}`).
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal",
             *       "targetServices": [
             *         "gateway",
             *         "hot-deals",
             *         "fashion-bestseller"
             *       ],
             *       "httpEndpointZones": {
             *         "type": "select",
             *         "targetType": "com.steadybit.extension_container.container",
             *         "attribute": "aws.zone",
             *         "mode": "fixed",
             *         "count": 2
             *       }
             *     }
             */
            experimentVariables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * @description An optional external identifier used for create-or-update semantics.
             * @example 1234567
             */
            externalId?: string;
            /**
             * @deprecated
             * @description An optional external reference. Will be removed and is replaced by tags. If used with experiment creation, the value will be added as a tag.
             * @example INCIDENT-4711
             */
            externalReference?: string;
            /**
             * @description The hypothesis that is validated by the experiment
             * @example System is able to survive a latency in the network of 1500ms
             */
            hypothesis?: string;
            /**
             * @description The lanes (steps executed in parallel) in the experiment. Each lane consists of multiple steps that are executed sequential per lane.
             * @example [
             *       {
             *         "steps": [
             *           {
             *             "type": "action",
             *             "ignoreFailure": false,
             *             "parameters": {
             *               "duration": "30s"
             *             },
             *             "actionType": "com.steadybit.extension_host.stress-cpu",
             *             "radius": {
             *               "targetType": "com.steadybit.extension_container.container",
             *               "predicate": {
             *                 "operator": "AND",
             *                 "predicates": [
             *                   {
             *                     "key": "k8s.deployment",
             *                     "operator": "EQUALS",
             *                     "values": [
             *                       "hot-deals"
             *                     ]
             *                   }
             *                 ]
             *               },
             *               "query": null,
             *               "percentage": 100
             *             }
             *           }
             *         ]
             *       }
             *     ]
             */
            lanes: components["schemas"]["ExperimentLaneAO"][];
            /**
             * @description Name of the experiment to easily identify the experiment
             * @example Shop survives unavailability of hot-deals products
             */
            name: string;
            /**
             * @description The properties of the experiment
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * @description Team keys with which the experiment is shared with
             * @example [OPS, SHOP]
             */
            sharedTeams?: string[];
            /**
             * @description An optional set of tags you can use to search for.
             * @example [
             *       "myTag",
             *       "myOtherTag"
             *     ]
             */
            tags?: string[];
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            team: string;
        };
        /**
         * @description Create or update an experiment based on an experiment template.
         * @example {
         *       "environment": "steadybit-demo",
         *       "team": "DEMO",
         *       "placeholders": [
         *         {
         *           "key": "CLUSTER",
         *           "value": "demo-cluster"
         *         },
         *         {
         *           "key": "BOOL",
         *           "value": true
         *         },
         *         {
         *           "key": "NUMBER",
         *           "value": 15
         *         },
         *         {
         *           "key": "KEYVALUE",
         *           "value": [
         *             {
         *               "key": "example-a",
         *               "value": "abc"
         *             },
         *             {
         *               "key": "example-b",
         *               "value": "123"
         *             }
         *           ]
         *         },
         *         {
         *           "key": "LIST",
         *           "value": [
         *             "entry1",
         *             "entry2",
         *             "entry3"
         *           ]
         *         },
         *         {
         *           "key": "FILE",
         *           "value": {
         *             "fileName": "example.txt",
         *             "data": "SGVsbG8gV29ybGQh"
         *           }
         *         }
         *       ],
         *       "externalId": "1234567",
         *       "experimentVariables": {
         *         "httpEndpoint": "http://dev.shop.products.internal"
         *       }
         *     }
         */
        CreateExperimentFromTemplateAO: {
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment?: string;
            /**
             * @description Variables that will be added to the created experiment design. A `key` that exists already in the environment variables will be overridden. Each value is either a constant string, an array of constant strings, or a select expression object.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal"
             *     }
             */
            experimentVariables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * @description An optional external identifier used for create-or-update semantics.
             * @example 1234567
             */
            externalId?: string;
            /** @description List of template placeholder values */
            placeholders?: components["schemas"]["ExperimentTemplatePlaceholderValueAO"][];
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            team: string;
        };
        CursorSliceResponseAOTargetAO: {
            /**
             * @description Are there more items that can be fetched with the given nextCursor?
             * @example true
             */
            hasNext?: boolean;
            items?: components["schemas"]["TargetAO"][];
            /**
             * @description The cursor to use to fetch the next page
             * @example eyJhZ2VudElkIjogImFnZW50LTEyMyIsICJuYW1lIjogImRlcGxveW1lbnQtYSIsICJ0eXBlIjogImNvbS5zdGVhZHliaXQuZXh0ZW5zaW9uX2t1YmVybmV0ZXMua3ViZXJuZXRlcy1kZXBsb3ltZW50In0=
             */
            nextCursor?: string | null;
        };
        CustomWebhookAO: {
            /**
             * @description The events that are being sent or a list containing a single `*` if all supported event types should be used.
             *
             *     Supported Events:
             *      - "experiment.execution.requested"
             *      - "experiment.execution.created"
             *      - "experiment.execution.preflight"
             *      - "experiment.execution.completed"
             *      - "experiment.execution.failed"
             *      - "experiment.execution.errored"
             *      - "experiment.execution.canceled"
             *      - "experiment.execution.step-started"
             *      - "experiment.execution.step-completed"
             *      - "experiment.execution.step-failed"
             *      - "experiment.execution.step-errored"
             *      - "experiment.execution.step-canceled"
             *      - "experiment.execution.step-skipped"
             *      - "killswitch.engaged"
             *      - "killswitch.disengaged"
             * @example [
             *       "experiment.execution.created",
             *       "experiment.execution.completed"
             *     ]
             */
            events: string[];
            /**
             * @description Additional headers to include in the webhook request.
             * @example {
             *       "X-Custom-Header": "CustomValue",
             *       "X-Another-Header": "AnotherValue"
             *     }
             */
            headers?: {
                [key: string]: string;
            };
            /**
             * Format: uuid
             * @description The id of the webhook
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id: string;
            /**
             * @description The name of the webhook
             * @example Custom Webhook
             */
            name: string;
            /**
             * @description The scope of the webhook / integration
             * @example TEAM
             * @enum {string}
             */
            scope: "GLOBAL" | "TEAM";
            /**
             * @description If a secret is provided a signature of the body is computed using `HMAC SHA-256` and sent as `X-SB-Signature` http header. You can use this header to verify the message.
             * @example secret123!!
             */
            secret?: string;
            /**
             * @description The body size can get very large as we include all target attributes for each target of your experiments. When having experiments with many targets, it might be useful to filter the attributes to only include the ones you are interested in. You can use the wildcard character '*' to match all attributes or a comma-separated-list of attribute-names. If the field is empty, no attributes will be included.
             * @example [
             *       "k8s.cluster-name",
             *       "k8s.deployment"
             *     ]
             */
            targetAttributeIncludes: string[];
            /**
             * @description The key of the team if the scope is `TEAM`
             * @example ADM
             */
            team?: string;
            /**
             * @description The URL of the webhook
             * @example https://example.com/webhook
             */
            url: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version: number;
        };
        CustomWebhookUpsertAO: {
            /**
             * @description The events that are being sent or a list containing a single `*` if all supported event types should be used.
             *
             *     Supported Events:
             *      - "experiment.execution.requested"
             *      - "experiment.execution.created"
             *      - "experiment.execution.preflight"
             *      - "experiment.execution.completed"
             *      - "experiment.execution.failed"
             *      - "experiment.execution.errored"
             *      - "experiment.execution.canceled"
             *      - "experiment.execution.step-started"
             *      - "experiment.execution.step-completed"
             *      - "experiment.execution.step-failed"
             *      - "experiment.execution.step-errored"
             *      - "experiment.execution.step-canceled"
             *      - "experiment.execution.step-skipped"
             *      - "killswitch.engaged"
             *      - "killswitch.disengaged"
             * @example [
             *       "experiment.execution.created",
             *       "experiment.execution.completed"
             *     ]
             */
            events: string[];
            /**
             * @description Additional headers to include in the webhook request.
             * @example {
             *       "X-Custom-Header": "CustomValue",
             *       "X-Another-Header": "AnotherValue"
             *     }
             */
            headers?: {
                [key: string]: string;
            };
            /**
             * Format: uuid
             * @description The id of the webhook or null if a new webhook should be created.
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id?: string | null;
            /**
             * @description The name of the webhook
             * @example Custom Webhook
             */
            name: string;
            /**
             * @description The scope of the webhook / integration
             * @example TEAM
             * @enum {string}
             */
            scope: "GLOBAL" | "TEAM";
            /**
             * @description If a secret is provided a signature of the body is computed using `HMAC SHA-256` and sent as `X-SB-Signature` http header. You can use this header to verify the message.
             * @example secret123!!
             */
            secret?: string;
            /**
             * @description The body size can get very large as we include all target attributes for each target of your experiments. When having experiments with many targets, it might be useful to filter the attributes to only include the ones you are interested in. You can use the wildcard character '*' to match all attributes or a comma-separated-list of attribute-names. If the field is empty, no attributes will be included.
             * @example [
             *       "k8s.cluster-name",
             *       "k8s.deployment"
             *     ]
             */
            targetAttributeIncludes: string[];
            /**
             * @description The key of the team if the scope is `TEAM`
             * @example ADM
             */
            team?: string;
            /**
             * @description The URL of the webhook
             * @example https://example.com/webhook
             */
            url: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /** @description Default blast radius configuration for an action. */
        DefaultBlastRadiusAO: {
            /**
             * @description Mode of the default blast radius.
             * @example PERCENTAGE
             */
            mode: string;
            /**
             * Format: int32
             * @description Value for the mode. Percentage (0-100) when mode is PERCENTAGE, max target count when mode is MAXIMUM.
             * @example 100
             */
            value?: number;
        };
        /**
         * @description The environment in which the event was triggered
         * @example {
         *       "id": "1avfd231-8322-42f2-bad9-307dc962ec37",
         *       "name": "Global",
         *       "predicate": {
         *         "operator": "AND",
         *         "predicates": []
         *       }
         *     }
         */
        EnvironmentAL: {
            /** Format: uuid */
            id: string;
            name: string;
            predicate: components["schemas"]["TargetPredicateAO"];
        } | null;
        /**
         * @description An environment for limiting the access to discovered systems for a team.
         * @example {
         *       "id": "2v1av42-e525-4c00-a13a-1ac32d170724",
         *       "name": "Global",
         *       "version": 0,
         *       "query": "aws.account=\"123\" OR aws.account=\"456\"",
         *       "state": "READY"
         *     }
         */
        EnvironmentAO: {
            /**
             * Format: uuid
             * @description Unique identifier of a environment
             */
            id?: string;
            /**
             * @description Name of the environment.
             * @example Global
             */
            name: string;
            predicate: components["schemas"]["TargetPredicateAO"];
            /**
             * @description Alternative to `predicate`. If both `query` and `predicate` will be provided, `query` will override the `predicate`.
             * @example (aws.account="123" OR aws.account="456"
             */
            query?: string | null;
            /**
             * @description State of the environment to indicate current background tasks.
             * @example "READY"
             * @enum {string}
             */
            state: "CREATED" | "UPDATED" | "READY" | "ERROR" | "UNKNOWN";
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number;
        };
        /**
         * @description List of environments.
         * @example {
         *       "environments": [
         *         {
         *           "id": "2v1av42-e525-4c00-a13a-1ac32d170724",
         *           "name": "Global",
         *           "version": 0,
         *           "query": "aws.account=\"123\" OR aws.account=\"456\"",
         *           "state": "READY"
         *         }
         *       ]
         *     }
         */
        EnvironmentSummariesAO: {
            environments?: components["schemas"]["EnvironmentAO"][];
        };
        /**
         * @description Experiment execution data that should be used only for that specific experiment execution and will not update the experiment design.
         * @example {}
         */
        ExecuteExperimentRequestAO: {
            /**
             * @description The name of the environment with which the experiment execution should be overridden once and executed in
             * @example Shop Stage
             */
            environment?: string;
            /**
             * @description Variables that will be merged for the single experiment execution with the variables defined of the experiment or environment that the experiment will be executed in. A `key` that exists already in the experiment or environment variables will be overridden for this execution, all others will be added solely in the context of this experiment execution. Existing variables don't have to be repeated in this parameter. Each value is either a constant string, an array of constant strings, or a select expression object.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal"
             *     }
             */
            variables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
        } | null;
        /**
         * @description A single experiment execution that was triggered from a single experiment.
         * @example {
         *       "id": 58070,
         *       "key": "SHOP-1"
         *       "apiLocation": "https://api.steadybit.com/experiments/execute/SHOP-1",
         *       "uiLocation": "https://platform.steadybit.com/experiments/edit/SHOP-1/executions/1234"
         *     }
         */
        ExecuteExperimentResponseAO: {
            /**
             * @description A link to the API for the experiment execution
             * @example https://api.steadybit.com/experiments/execute/SHOP-1
             */
            apiLocation: string;
            /**
             * Format: int64
             * @description Unique experiment execution id that identifies this specific experiment execution
             * @example 1234
             */
            executionId?: number;
            /**
             * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
             * @example SHOP-1
             */
            key: string;
            /**
             * @description A link to the UI for the experiment execution
             * @example https://platform.steadybit.com/experiments/edit/SHOP-1/executions/1234
             */
            uiLocation: string;
        };
        /**
         * @description Modifications that should be applied to the execution if the preflight check was successful.
         * @example [
         *       {
         *         "type": "set_property_value",
         *         "propertyKey": "approvedBy",
         *         "value": "Daniel"
         *       },
         *       {
         *         "type": "add_value_to_list_property",
         *         "propertyKey": "observations",
         *         "value": "This looks interesting!"
         *       }
         *     ]
         */
        ExecutionModification: {
            type: string;
        };
        /**
         * @example {
         *       "key": "ADM-2",
         *       "name": "Blackhole Hot-deals",
         *       "team": "ADM",
         *       "sharedTeams": [
         *         "SHOP"
         *       ],
         *       "environment": "Global",
         *       "created": "2023-05-03T08:24:30.183237Z",
         *       "createdBy": {
         *         "username": "ag1hb7ap-d299-47ab-998f-c2a53b433820",
         *         "name": "Manuel",
         *         "pictureUrl": "https://.../picture.png"
         *       },
         *       "edited": "2023-05-03T13:31:12.533664Z",
         *       "editedBy": {
         *         "username": "ag1hb7ap-d299-47ab-998f-c2a53b433820",
         *         "name": "Manuel",
         *         "pictureUrl": "https://.../picture.png"
         *       },
         *       "lanes": [
         *         {
         *           "steps": [
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "30s"
         *               },
         *               "actionType": "com.steadybit.extension_host.stress-cpu",
         *               "radius": {
         *                 "targetType": "com.steadybit.extension_container.container",
         *                 "predicate": {
         *                   "operator": "AND",
         *                   "predicates": [
         *                     {
         *                       "key": "k8s.deployment",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "hot-deals"
         *                       ]
         *                     }
         *                   ]
         *                 },
         *                 "query": null,
         *                 "percentage": 100
         *               }
         *             }
         *           ]
         *         }
         *       ],
         *       "properties": {
         *         "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
         *       }
         *     }
         */
        ExperimentAO: {
            /**
             * Format: date-time
             * @description Timestamp when the experiment was created
             * @example 2023-01-01T09:00:00Z
             */
            created: string;
            createdBy: components["schemas"]["UserSummaryAO"];
            /**
             * Format: date-time
             * @description Timestamp when the experiment was edited the last time
             * @example 2023-01-01T09:00:00Z
             */
            edited: string;
            editedBy: components["schemas"]["UserSummaryAO"];
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment?: string;
            /**
             * @description Variables that will be used when the experiment will be executed. Experiment variables will override existing environment variables. Each value is either a constant string, an array of constant strings, or a select expression object (`{"type":"select",...}`).
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal",
             *       "targetServices": [
             *         "gateway",
             *         "hot-deals",
             *         "fashion-bestseller"
             *       ],
             *       "httpEndpointZones": {
             *         "type": "select",
             *         "targetType": "com.steadybit.extension_container.container",
             *         "attribute": "aws.zone",
             *         "mode": "fixed",
             *         "count": 2
             *       }
             *     }
             */
            experimentVariables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * @description An optional external identifier used for create-or-update semantics.
             * @example 1234567
             */
            externalId?: string;
            /**
             * @deprecated
             * @description An optional external reference. Will be removed and is replaced by tags. If used with experiment creation, the value will be added as a tag.
             * @example INCIDENT-4711
             */
            externalReference?: string;
            /**
             * @description The hypothesis that is validated by the experiment
             * @example System is able to survive a latency in the network of 1500ms
             */
            hypothesis?: string;
            /**
             * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
             * @example ADM-2
             */
            key: string;
            /**
             * @description The lanes (steps executed in parallel) in the experiment. Each lane consists of multiple steps that are executed sequential per lane.
             * @example [
             *       {
             *         "steps": [
             *           {
             *             "type": "action",
             *             "ignoreFailure": false,
             *             "parameters": {
             *               "duration": "30s"
             *             },
             *             "actionType": "com.steadybit.extension_host.stress-cpu",
             *             "radius": {
             *               "targetType": "com.steadybit.extension_container.container",
             *               "predicate": {
             *                 "operator": "AND",
             *                 "predicates": [
             *                   {
             *                     "key": "k8s.deployment",
             *                     "operator": "EQUALS",
             *                     "values": [
             *                       "hot-deals"
             *                     ]
             *                   }
             *                 ]
             *               },
             *               "query": null,
             *               "percentage": 100
             *             }
             *           }
             *         ]
             *       }
             *     ]
             */
            lanes: components["schemas"]["ExperimentLaneAO"][];
            /**
             * @description Name of the experiment to easily identify the experiment
             * @example Shop survives unavailability of hot-deals products
             */
            name: string;
            /**
             * @description The properties of the experiment
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * @description Team keys with which the experiment is shared with
             * @example [OPS, SHOP]
             */
            sharedTeams?: string[];
            /**
             * @description An optional set of tags you can use to search for.
             * @example [
             *       "myTag",
             *       "myOtherTag"
             *     ]
             */
            tags?: string[];
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            team: string;
            /** @description The placeholders that were used to create this experiment from a template */
            templatePlaceholders?: components["schemas"]["ExperimentTemplatePlaceholderValueAO"][];
            /** @description The title of the template that was used to create this experiment */
            templateTitle?: string;
            /**
             * Format: int32
             * @description Experiment database version.
             * @example 1
             */
            version?: number;
        };
        /**
         * @description A single experiment execution that was triggered from a single experiment.
         * @example {
         *       "id": 58070,
         *       "key": "SHOP-1",
         *       "name": "Shop should survive a single pod outage",
         *       "hypothesis": "When a single container from steadybit-demo/fashion-bestseller fails the shop is still working as expected.",
         *       "requested": "2023-01-01T09:00:00.000000Z",
         *       "created": "2023-01-01T09:00:01.000000Z",
         *       "createdBy": {
         *         "username": "ag1hb7ap-d299-47ab-998f-c2a53b433820",
         *         "name": "Manuel",
         *         "pictureUrl": "https://.../picture.png"
         *       },
         *       "createdVia": "UI",
         *       "experimentVersion": "5",
         *       "ended": "2023-01-01T09:10:00.000000Z",
         *       "state": "FAILED",
         *       "reason": "Check failure."
         *     }
         */
        ExperimentExecutionAO: {
            canceledBy?: components["schemas"]["UserSummaryAO"];
            /**
             * Format: date-time
             * @description Timestamp when the experiment was created
             * @example 2023-01-01T09:00:01Z
             */
            created?: string;
            createdBy?: components["schemas"]["UserSummaryAO"];
            /**
             * @description The creation trigger that caused this experiment execution to be started
             * @example UI
             * @enum {string}
             */
            createdVia?: "API" | "CLI" | "UI" | "SCHEDULE" | "SUITE" | "MCP";
            /**
             * Format: date-time
             * @description Timestamp when the experiment ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * Format: int32
             * @description Experiment design version which can be used to identify changes between experiment runs
             * @example 5
             */
            experimentVersion?: number;
            /**
             * @description The hypothesis that is validated by the experiment
             * @example System is able to survive a latency in the network of 1500ms
             */
            hypothesis?: string;
            /**
             * Format: int32
             * @description Unique experiment execution id that identifies this specific experiment execution
             * @example 1523
             */
            id?: number;
            /**
             * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
             * @example ADM-2
             */
            key?: string;
            /**
             * @description Name of the experiment to easily identify the experiment
             * @example Shop survives unavailability of hot-deals products
             */
            name?: string;
            /**
             * @description The properties of the experiment execution
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "Chuck Norris allows that execution"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * Format: int32
             * @description Version of the properties for optimistic locking (optional in the Update-API)
             * @example 1
             */
            propertiesVersion?: number;
            /**
             * @description Reason in case the experiment execution failed or errored
             * @example Action error
             */
            reason?: string;
            /**
             * Format: date-time
             * @description Timestamp when the experiment was requested
             * @example 2023-01-01T09:00:00Z
             */
            requested?: string;
            /**
             * Format: date-time
             * @description Timestamp when the experiment was started
             * @example 2023-01-01T09:00:02Z
             */
            started?: string;
            /**
             * @description Current state of the experiment (e.g. CREATED, RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
            /**
             * @description The steps that are executed in parallel or sequence in the experiment.
             * @example [
             *       {
             *         "ignoreFailure": false,
             *         "parameters": {
             *           "duration": "10s"
             *         }
             *       },
             *       {
             *         "predecessorId": "40b0f797-912d-4256-8887-1553561962a9",
             *         "ignoreFailure": false,
             *         "parameters": {
             *           "cpuLoad": 100,
             *           "workers": 0,
             *           "duration": "30s"
             *         },
             *         "actionId": "com.steadybit.extension_container.stress_cpu",
             *         "actionKind": "ATTACK",
             *         "radius": {
             *           "targetType": "com.steadybit.extension_container.container",
             *           "percentage": 50,
             *           "predicate": {
             *             "operator": "AND",
             *             "predicates": [
             *               {
             *                 "key": "container.host/name",
             *                 "operator": "EQUALS",
             *                 "values": [
             *                   "docker-desktop/minikube"
             *                 ]
             *               }
             *             ]
             *           }
             *         },
             *         "targetExecutions": [
             *           {
             *             "type": "com.steadybit.extension_container.container",
             *             "name": "docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea",
             *             "state": "COMPLETED",
             *             "attributes": [
             *               {
             *                 "key": "container.port",
             *                 "value": "51152:2376"
             *               },
             *               {
             *                 "key": "container.engine",
             *                 "value": "docker"
             *               },
             *               {
             *                 "key": "container.host/name",
             *                 "value": "docker-desktop/minikube"
             *               },
             *               {
             *                 "key": "container.host",
             *                 "value": "docker-desktop"
             *               }
             *             ]
             *           }
             *         ],
             *         "totalTargetCount": 1
             *       }
             *     ]
             */
            steps?: components["schemas"]["AbstractExperimentExecutionStepAO"][];
            /**
             * @description Tags of the experiment at the time the execution was requested
             * @example [
             *       "resilience",
             *       "shop"
             *     ]
             */
            tags?: string[];
            /**
             * @description Variables and their origins that have been used for this execution
             * @example {
             *       "httpEndpoint": {
             *         "value": "http://dev.shop.products.internal",
             *         "origin": "ENVIRONMENT"
             *       }
             *     }
             */
            variables?: {
                [key: string]: components["schemas"]["ExperimentExecutionVariableAO"];
            };
        };
        ExperimentExecutionPageItemAO: {
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution was created
             * @example 2023-01-01T09:00:01Z
             */
            created?: string;
            createdBy?: string;
            createdByDetails?: components["schemas"]["UserSummaryAO"];
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution was ended
             * @example 2023-01-01T09:00:01Z
             */
            ended?: string;
            /**
             * @description The name of the environment that this execution was using
             * @example Global
             */
            environment?: string;
            /**
             * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
             * @example ADM-2
             */
            experimentKey?: string;
            /**
             * Format: int64
             * @description Unique experiment execution id that identifies a single experiment execution
             * @example 123
             */
            id?: number;
            /**
             * @description Name of the experiment to easily identify the experiment
             * @example Shop survives unavailability of hot-deals products
             */
            name?: string;
            /**
             * @description The properties of the experiment execution
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "Chuck Norris allows that execution"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /** @description Details about the failure/error reason. */
            reason?: string;
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution was requested
             * @example 2023-01-01T09:00:01Z
             */
            requested?: string;
            /**
             * @description Was this execution triggered by a schedule?
             * @example true
             */
            scheduled?: boolean;
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution was started
             * @example 2023-01-01T09:00:01Z
             */
            started?: string;
            /**
             * @description Current state of the experiment execution (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
            /**
             * @description The key of the team that this experiment is assigned to
             * @example ADM
             */
            teamKey?: string;
        };
        /** @description Filter for experiment execution report data, optionally scoped to specific teams, environments, and services. */
        ExperimentExecutionReportFilterAO: {
            /** @description Restrict results to the given environment IDs. If not provided, all environments are included. */
            environmentIds?: string[] | null;
            /**
             * Format: date
             * @description Start date of the report range (inclusive).
             * @example 2026-01-01
             */
            from: string;
            /**
             * @description The time bucket granularity for report aggregation.
             * @example MONTHLY
             * @enum {string}
             */
            rollup?: "MONTHLY" | "DAILY";
            /** @description Restrict results to the given service IDs. If not provided, all services are included. */
            serviceIds?: string[] | null;
            /** @description Restrict results to the given team IDs. If not provided, all teams are included. */
            teamIds?: string[] | null;
            /**
             * Format: date
             * @description End date of the report range (inclusive).
             * @example 2026-03-01
             */
            to: string;
        };
        /**
         * @description Filters are defined in the body of the request.
         * @example {
         *       "page": 0,
         *       "environments": [
         *         "Global"
         *       ],
         *       "experimentKeys": [
         *         "ADM-9"
         *       ],
         *       "teamKeys": [
         *         "GITHUB"
         *       ],
         *       "teamKeysExclude": [
         *         "ADM"
         *       ],
         *       "services": [
         *         "shopping-cart"
         *       ],
         *       "states": [
         *         "errored",
         *         "canceled"
         *       ],
         *       "requestedFrom": "2024-05-17T00:00:00Z",
         *       "requestedTo": "2024-06-24T00:00:00Z",
         *       "endedFrom": "2024-05-17T00:00:00Z",
         *       "endedTo": "2024-06-24T00:00:00Z"
         *     }
         */
        ExperimentExecutionsRequestAO: {
            /**
             * Format: date-time
             * @description Filter results by range of created date
             * @example 2021-01-01T00:00:00Z
             */
            createdFrom?: string | null;
            /**
             * Format: date-time
             * @description Filter results by range of created date
             * @example 2021-01-01T00:00:00Z
             */
            createdTo?: string | null;
            /**
             * Format: date-time
             * @description Filter results by range of ended date
             * @example 2021-01-01T00:00:00Z
             */
            endedFrom?: string | null;
            /**
             * Format: date-time
             * @description Filter results by range of ended date
             * @example 2021-01-01T00:00:00Z
             */
            endedTo?: string | null;
            /**
             * @description Filter results by one or more environments
             * @example [
             *       "Global"
             *     ]
             */
            environments?: string[] | null;
            /**
             * @description Filter results by one or more experiment-keys
             * @example [
             *       "ADM-9"
             *     ]
             */
            experimentKeys?: string[] | null;
            /**
             * @description Filter results by name and/or key of the experiment
             * @example Outage
             */
            name?: string | null;
            /** Format: int32 */
            page?: number;
            /**
             * Format: date-time
             * @description Filter results by range of requested date
             * @example 2021-01-01T00:00:00Z
             */
            requestedFrom?: string | null;
            /**
             * Format: date-time
             * @description Filter results by range of requested date
             * @example 2021-01-01T00:00:00Z
             */
            requestedTo?: string | null;
            /**
             * @description Filter results by one or more service names that should be included in the result
             * @example [
             *       "shopping-cart"
             *     ]
             */
            services?: string[] | null;
            /** Format: int32 */
            size?: number;
            /**
             * @description Filter results by one or more states. Possible values: [CREATED, PREPARED, RUNNING, FAILED, CANCELED, COMPLETED, ERRORED]
             * @example [
             *       "CREATED"
             *     ]
             */
            states?: string[] | null;
            /**
             * @description Filter results by one or more team-keys that should be included in the result
             * @example [
             *       "ADM"
             *     ]
             */
            teamKeys?: string[] | null;
            /**
             * @description Filter results by one or more team-keys that should be excluded in the result
             * @example [
             *       "ADM"
             *     ]
             */
            teamKeysExclude?: string[] | null;
        };
        /**
         * @description An action-step that is executed as part of an experiment.
         * @example {
         *       "stepType": "ACTION",
         *       "id": "0199c3f2-48c7-706d-b102-9cb09dd41b5d",
         *       "state": "COMPLETED",
         *       "started": "2025-10-08T13:11:01.487541Z",
         *       "ended": "2025-10-08T13:11:11.490268Z",
         *       "predecessorId": "40b0f797-912d-4256-8887-1553561962a9",
         *       "ignoreFailure": false,
         *       "parameters": {
         *         "cpuLoad": 100,
         *         "workers": 0,
         *         "duration": "30s"
         *       },
         *       "actionId": "com.steadybit.extension_container.stress_cpu",
         *       "actionKind": "ATTACK",
         *       "radius": {
         *         "targetType": "com.steadybit.extension_container.container",
         *         "percentage": 50,
         *         "predicate": {
         *           "operator": "AND",
         *           "predicates": [
         *             {
         *               "key": "container.host/name",
         *               "operator": "EQUALS",
         *               "values": [
         *                 "docker-desktop/minikube"
         *               ]
         *             }
         *           ]
         *         }
         *       },
         *       "targetExecutions": [
         *         {
         *           "type": "com.steadybit.extension_container.container",
         *           "name": "docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea",
         *           "state": "COMPLETED",
         *           "attributes": [
         *             {
         *               "key": "container.port",
         *               "value": "51152:2376"
         *             },
         *             {
         *               "key": "container.engine",
         *               "value": "docker"
         *             },
         *             {
         *               "key": "container.host/name",
         *               "value": "docker-desktop/minikube"
         *             },
         *             {
         *               "key": "container.host",
         *               "value": "docker-desktop"
         *             }
         *           ]
         *         }
         *       ],
         *       "totalTargetCount": 1
         *     }
         */
        ExperimentExecutionStepActionAO: {
            /**
             * @description Unique identifier of the action that is executed in this step
             * @example com.steadybit.extension_container.stress_cpu
             */
            actionId?: string;
            /**
             * @description Kind of the action (e.g. attack, check, loadtest)
             * @example ATTACK
             * @enum {string}
             */
            actionKind?: "ATTACK" | "CHECK" | "LOAD_TEST" | "OTHER" | "BASIC";
            /**
             * @description Custom label assigned during experiment design to express the intention of this step
             * @example Container 'xyz' can not be reached
             */
            customLabel?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * Format: uuid
             * @description Unique identifier of this step execution
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            id?: string;
            /**
             * @description Whether the experiment should fail/error immediately in case this step fails/errors.
             * @example false
             */
            ignoreFailure?: boolean;
            /**
             * @description Step-specific parameters of the experiment step configuration
             * @example {
             *       "duration": "10s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            /**
             * Format: uuid
             * @description Unique identifier of the step execution that precedes this step, null if it is the first step of a lane
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            predecessorId?: string;
            radius?: components["schemas"]["BlastRadiusAO"];
            /**
             * @description Reason in case this experiment step execution failed or errored
             * @example Couldn't read state of container...
             */
            reason?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step was started
             * @example 2023-01-01T09:00:00Z
             */
            started?: string;
            /**
             * @description Current state of this step in the experiment (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
            /**
             * @description Type of this step execution (e.g. ACTION, WAIT) (enum property replaced by openapi-typescript)
             * @enum {string}
             */
            stepType: "ACTION";
            /**
             * @description List of targets that are expected to be effected by this action. This list may change in case targets aren't available at the specific time of execution
             * @example [
             *       {
             *         "type": "com.steadybit.extension_container.container",
             *         "name": "docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea",
             *         "state": "COMPLETED",
             *         "attributes": [
             *           {
             *             "key": "container.port",
             *             "value": "51152:2376"
             *           },
             *           {
             *             "key": "container.engine",
             *             "value": "docker"
             *           },
             *           {
             *             "key": "container.host/name",
             *             "value": "docker-desktop/minikube"
             *           },
             *           {
             *             "key": "container.host",
             *             "value": "docker-desktop"
             *           }
             *         ]
             *       }
             *     ]
             */
            targetExecutions?: components["schemas"]["TargetExecutionAO"][];
            /**
             * Format: int64
             * @description Amount of targets that are effect int total
             * @example 23
             */
            totalTargetCount?: number;
        };
        /**
         * @description A service validation step that is executed as part of an experiment.
         * @example {
         *       "stepType": "SERVICE-VALIDATION",
         *       "id": "40b0f797-912d-4256-8887-1553561962a9",
         *       "state": "COMPLETED",
         *       "started": "2025-06-18T08:32:01.850479Z",
         *       "ended": "2025-06-18T08:32:11.886043Z",
         *       "predecessorId": null,
         *       "ignoreFailure": false,
         *       "parameters": {
         *         "duration": "60s"
         *       },
         *       "serviceId": "cc06f132-0694-4ffa-aee2-13d8fafa3a8b",
         *       "validations": []
         *     }
         */
        ExperimentExecutionStepServiceValidationAO: {
            /**
             * @description Custom label assigned during experiment design to express the intention of this step
             * @example Container 'xyz' can not be reached
             */
            customLabel?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * Format: uuid
             * @description Unique identifier of this step execution
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            id?: string;
            /**
             * @description Whether the experiment should fail/error immediately in case this step fails/errors.
             * @example false
             */
            ignoreFailure?: boolean;
            /**
             * @description Step-specific parameters of the experiment step configuration
             * @example {
             *       "duration": "10s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            /**
             * Format: uuid
             * @description Unique identifier of the step execution that precedes this step, null if it is the first step of a lane
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            predecessorId?: string;
            /**
             * @description Reason in case this experiment step execution failed or errored
             * @example Couldn't read state of container...
             */
            reason?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step was started
             * @example 2023-01-01T09:00:00Z
             */
            started?: string;
            /**
             * @description Current state of this step in the experiment (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
            /**
             * @description Type of this step execution (e.g. ACTION, WAIT)
             * @example ACTION
             */
            stepType: string;
        } & {
            /**
             * Format: uuid
             * @description Unique identifier of the service.
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            serviceId?: string;
            /** @description List of actions performed as part of this service validation step. */
            validations?: components["schemas"]["ExperimentExecutionStepActionAO"][];
        } & {
            /**
             * @description discriminator enum property added by openapi-typescript
             * @enum {string}
             */
            stepType: "SERVICE-VALIDATION";
        };
        /**
         * @description A wait step that is executed as part of an experiment.
         * @example {
         *       "stepType": "WAIT",
         *       "id": "40b0f797-912d-4256-8887-1553561962a9",
         *       "state": "COMPLETED",
         *       "started": "2025-06-18T08:32:01.850479Z",
         *       "ended": "2025-06-18T08:32:11.886043Z",
         *       "predecessorId": null,
         *       "ignoreFailure": false,
         *       "parameters": {
         *         "duration": "10s"
         *       }
         *     }
         */
        ExperimentExecutionStepWaitAO: {
            /**
             * @description Custom label assigned during experiment design to express the intention of this step
             * @example Container 'xyz' can not be reached
             */
            customLabel?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * Format: uuid
             * @description Unique identifier of this step execution
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            id?: string;
            /**
             * @description Whether the experiment should fail/error immediately in case this step fails/errors.
             * @example false
             */
            ignoreFailure?: boolean;
            /**
             * @description Step-specific parameters of the experiment step configuration
             * @example {
             *       "duration": "10s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            /**
             * Format: uuid
             * @description Unique identifier of the step execution that precedes this step, null if it is the first step of a lane
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            predecessorId?: string;
            /**
             * @description Reason in case this experiment step execution failed or errored
             * @example Couldn't read state of container...
             */
            reason?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step was started
             * @example 2023-01-01T09:00:00Z
             */
            started?: string;
            /**
             * @description Current state of this step in the experiment (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
            /**
             * @description Type of this step execution (e.g. ACTION, WAIT) (enum property replaced by openapi-typescript)
             * @enum {string}
             */
            stepType: "WAIT";
        };
        /**
         * @description List of experiment exeuctions.
         * @example {
         *       "executions": [
         *         {
         *           "id": 102,
         *           "key": "SHOP-1",
         *           "name": "Shop survives outage of a single pod",
         *           "created": "2023-01-01T09:00:01.000000Z",
         *           "ended": "2023-01-01T09:01:00.000000Z",
         *           "state": "FAILED"
         *         },
         *         {
         *           "id": 103,
         *           "key": "SHOP-1",
         *           "name": "Shop survives outage of a single pod",
         *           "created": "2023-01-01T09:10:00.000000Z",
         *           "ended": "2023-01-01T09:11:00.000000Z",
         *           "state": "COMPLETED"
         *         },
         *         {
         *           "id": 110,
         *           "key": "SHOP-2",
         *           "name": "DataDog monitors notices pod unavailability",
         *           "created": "2023-01-01T09:30:00.000000Z",
         *           "ended": "2023-01-01T09:41:00.000000Z",
         *           "state": "COMPLETED"
         *         }
         *       ]
         *     }
         */
        ExperimentExecutionSummariesAO: {
            /**
             * @description List of experiment executions
             * @example [
             *       {
             *         "id": 102,
             *         "key": "SHOP-1",
             *         "name": "Shop survives outage of a single pod",
             *         "requested": "2023-01-01T09:00:00.000000Z",
             *         "created": "2023-01-01T09:00:01.000000Z",
             *         "started": "2023-01-01T09:00:02.000000Z",
             *         "ended": "2023-01-01T09:01:00.000000Z",
             *         "state": "FAILED"
             *       },
             *       {
             *         "id": 103,
             *         "key": "SHOP-1",
             *         "name": "Shop survives outage of a single pod",
             *         "requested": "2023-01-01T09:10:00.000000Z",
             *         "created": "2023-01-01T09:10:01.000000Z",
             *         "started": "2023-01-01T09:10:02.000000Z",
             *         "ended": "2023-01-01T09:11:00.000000Z",
             *         "state": "COMPLETED"
             *       },
             *       {
             *         "id": 110,
             *         "key": "SHOP-2",
             *         "name": "DataDog monitors notices pod unavailability",
             *         "requested": "2023-01-01T09:30:00.000000Z",
             *         "created": "2023-01-01T09:30:01.000000Z",
             *         "started": "2023-01-01T09:30:02.000000Z",
             *         "ended": "2023-01-01T09:41:00.000000Z",
             *         "state": "COMPLETED"
             *       }
             *     ]
             */
            executions?: components["schemas"]["ExperimentExecutionSummaryAO"][];
        };
        /**
         * @description List of experiment executions
         * @example [
         *       {
         *         "id": 102,
         *         "key": "SHOP-1",
         *         "name": "Shop survives outage of a single pod",
         *         "requested": "2023-01-01T09:00:00.000000Z",
         *         "created": "2023-01-01T09:00:01.000000Z",
         *         "started": "2023-01-01T09:00:02.000000Z",
         *         "ended": "2023-01-01T09:01:00.000000Z",
         *         "state": "FAILED"
         *       },
         *       {
         *         "id": 103,
         *         "key": "SHOP-1",
         *         "name": "Shop survives outage of a single pod",
         *         "requested": "2023-01-01T09:10:00.000000Z",
         *         "created": "2023-01-01T09:10:01.000000Z",
         *         "started": "2023-01-01T09:10:02.000000Z",
         *         "ended": "2023-01-01T09:11:00.000000Z",
         *         "state": "COMPLETED"
         *       },
         *       {
         *         "id": 110,
         *         "key": "SHOP-2",
         *         "name": "DataDog monitors notices pod unavailability",
         *         "requested": "2023-01-01T09:30:00.000000Z",
         *         "created": "2023-01-01T09:30:01.000000Z",
         *         "started": "2023-01-01T09:30:02.000000Z",
         *         "ended": "2023-01-01T09:41:00.000000Z",
         *         "state": "COMPLETED"
         *       }
         *     ]
         */
        ExperimentExecutionSummaryAO: {
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution was created
             * @example 2023-01-01T09:00:01Z
             */
            created?: string;
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution ended
             * @example 2023-01-01T09:01:00Z
             */
            ended?: string;
            /**
             * Format: int32
             * @description Unique experiment execution id that identifies a single experiment execution
             * @example 123
             */
            id?: number;
            /**
             * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
             * @example ADM-2
             */
            key?: string;
            /**
             * @description Name of the experiment to easily identify the experiment
             * @example Shop survives unavailability of hot-deals products
             */
            name?: string;
            /**
             * @description The properties of the experiment execution
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "Chuck Norris allows that execution"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution was requested
             * @example 2023-01-01T09:00:00Z
             */
            requested?: string;
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution started
             * @example 2023-01-01T09:00:02Z
             */
            started?: string;
            /**
             * @description Current state of the experiment execution (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
        };
        /**
         * @description The variables resolved for this specific execution, keyed by name. Each entry carries the resolved value(s) and the tier the winning value originated from (ENVIRONMENT, SERVICE, EXPERIMENT, SCHEDULE, EXECUTION). A single-value variable's value is a string, a multi-value variable's value is an array of strings. Empty until the execution starts, as dynamic values are resolved once at run start and then stay stable for the whole run.
         * @example {
         *       "httpEndpoint": {
         *         "value": "http://shop.products.internal",
         *         "origin": "EXECUTION"
         *       }
         *     }
         */
        ExperimentExecutionVariableAO: {
            /** @enum {string} */
            origin?: "ENVIRONMENT" | "SERVICE" | "EXPERIMENT" | "SCHEDULE" | "EXECUTION";
            /** @description Either a single value (the common case, including single-element select results) or an array of values (multi-value select expressions). */
            value?: string | string[];
        };
        /**
         * @description A single lane of an experiment design. This lane can contain multiple steps that are executed sequentially
         * @example {
         *       "steps": [
         *         {
         *           "type": "action",
         *           "ignoreFailure": false,
         *           "parameters": {
         *             "duration": "30s"
         *           },
         *           "actionType": "com.steadybit.extension_host.stress-cpu",
         *           "radius": {
         *             "targetType": "com.steadybit.extension_container.container",
         *             "predicate": {
         *               "operator": "AND",
         *               "predicates": [
         *                 {
         *                   "key": "k8s.deployment",
         *                   "operator": "EQUALS",
         *                   "values": [
         *                     "hot-deals"
         *                   ]
         *                 }
         *               ]
         *             },
         *             "query": null,
         *             "percentage": 100
         *           }
         *         }
         *       ]
         *     }
         */
        ExperimentLaneAO: {
            /**
             * @description A list of steps that are executed sequentially in this lane.
             * @example [
             *       {
             *         "type": "action",
             *         "ignoreFailure": false,
             *         "parameters": {
             *           "duration": "30s"
             *         },
             *         "actionType": "com.steadybit.extension_host.stress-cpu",
             *         "radius": {
             *           "targetType": "com.steadybit.extension_container.container",
             *           "predicate": {
             *             "operator": "AND",
             *             "predicates": [
             *               {
             *                 "key": "k8s.deployment",
             *                 "operator": "EQUALS",
             *                 "values": [
             *                   "hot-deals"
             *                 ]
             *               }
             *             ]
             *           },
             *           "query": null,
             *           "percentage": 100
             *         }
             *       }
             *     ]
             */
            steps: components["schemas"]["BaseExperimentStepAO"][];
        };
        /** @description Filter for experiment report data, optionally scoped to specific teams and environments. */
        ExperimentReportFilterAO: {
            /** @description Restrict results to the given environment IDs. If not provided, all environments are included. */
            environmentIds?: string[] | null;
            /**
             * Format: date
             * @description Start date of the report range (inclusive).
             * @example 2026-01-01
             */
            from: string;
            /**
             * @description The time bucket granularity for report aggregation.
             * @example MONTHLY
             * @enum {string}
             */
            rollup?: "MONTHLY" | "DAILY";
            /** @description Restrict results to the given team IDs. If not provided, all teams are included. */
            teamIds?: string[] | null;
            /**
             * Format: date
             * @description End date of the report range (inclusive).
             * @example 2026-03-01
             */
            to: string;
        };
        /**
         * @description The risk for a single experiment linked to a service
         * @example {
         *       "experimentKey": "ADM-8",
         *       "risk": 42
         *     }
         */
        ExperimentRiskAO: {
            /** @description The experiment key */
            experimentKey?: string;
            /**
             * Format: int32
             * @description The calculated risk score for this experiment (0-100)
             */
            risk?: number;
        };
        /**
         * @description A schedule for an experiment.
         * @example {
         *       "experimentKey": "ADM-8",
         *       "cron": "30 * * * * ? *",
         *       "enabled": true,
         *       "allowParallel": true,
         *       "timezone": "Europe/Berlin",
         *       "variables": {},
         *       "id": "01951394-727f-76a0-8675-c7519ebd0ff5",
         *       "lastUpdated": "2025-02-17T11:04:10.623486Z",
         *       "editedBy": {
         *         "username": "ag1hb7ap-d299-47ab-998f-c2a53b433820",
         *         "name": "Manuel",
         *         "pictureUrl": "https://.../picture.png",
         *         "email": "manuel@example.org"
         *       },
         *       "nextExecution": "2025-02-20T05:54:30Z"
         *     }
         */
        ExperimentScheduleAO: {
            /**
             * @description Should the experiment run if another experiment is running? Default is true.
             * @example true
             */
            allowParallel?: boolean;
            /**
             * @description Cron expression for the experiment schedule. Can't be used in combination with `startAt`.
             * @example 0 15 10 ? * *
             */
            cron?: string | null;
            editedBy: components["schemas"]["UserSummaryAO"];
            /**
             * @description If `false`, the schedule is deactivated and no experiment will be executed. Default is true.
             * @example false
             */
            enabled?: boolean;
            /**
             * @description The experiment that should be scheduled.
             * @example ADM-123
             */
            experimentKey: string;
            id: string;
            /** Format: date-time */
            lastUpdated: string;
            /** Format: date-time */
            nextExecution?: string | null;
            /**
             * Format: date-time
             * @description Start date for a single execution. Can't be used in combination with `cron`.
             */
            startAt?: string | null;
            /**
             * @description Optional timezone for a experiment schedule. Can only be used with `cron`.
             * @example Europe/Berlin
             */
            timezone?: string | null;
            /**
             * @description Variables that will be used when the experiment will be executed. The variables will override existing environment or experiment variables. Each value is either a constant string, an array of constant strings, or a select expression object.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal"
             *     }
             */
            variables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
        };
        /**
         * @description A single step in a lane executing always exactly one action.
         * @example {
         *       "type": "action",
         *       "ignoreFailure": false,
         *       "parameters": {
         *         "duration": "30s"
         *       },
         *       "actionType": "com.steadybit.extension_host.stress-cpu",
         *       "radius": {
         *         "targetType": "com.steadybit.extension_container.container",
         *         "predicate": {
         *           "operator": "AND",
         *           "predicates": [
         *             {
         *               "key": "k8s.deployment",
         *               "operator": "EQUALS",
         *               "values": [
         *                 "hot-deals"
         *               ]
         *             }
         *           ]
         *         },
         *         "query": null,
         *         "percentage": 100
         *       }
         *     }
         */
        ExperimentStepActionAO: {
            customLabel?: string;
            /**
             * @description Ignore any errors and failures of this single step and continue the execution of an experiment run
             * @example false
             */
            ignoreFailure?: boolean;
            /** @description Optional metric checks used to define success or failure of this step */
            metricChecks?: components["schemas"]["MetricCheckAO"][];
            /** @description Optional metric queries used of this step to filter e.g. monitoring data */
            metricQueries?: components["schemas"]["MetricQueryAO"][];
            /**
             * @description Configuration parameters of this step that are saved during experiment design and evaluated at execution time.
             * @example {
             *       "duration": "30s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            type: string;
        } & {
            /**
             * @description The specific action that is used in this step
             * @example com.steadybit.extension_host.stress-cpu
             */
            actionType: string;
            radius?: components["schemas"]["ExperimentStepRadiusAO"];
        } & {
            /**
             * @description discriminator enum property added by openapi-typescript
             * @enum {string}
             */
            type: "action";
        };
        /**
         * @description Specifying the targets and random blast radius of the available targets
         * @example {
         *       "targetType": "com.steadybit.extension_container.container",
         *       "predicate": {
         *         "operator": "AND",
         *         "predicates": [
         *           {
         *             "key": "k8s.deployment",
         *             "operator": "EQUALS",
         *             "values": [
         *               "hot-deals"
         *             ]
         *           }
         *         ]
         *       }
         *     }
         */
        ExperimentStepRadiusAO: {
            advanced?: components["schemas"]["AdvancedRadiusAO"][];
            /** Format: int32 */
            maximum?: number;
            /** Format: int32 */
            percentage?: number;
            predicate?: components["schemas"]["TargetPredicateAO"];
            query?: components["schemas"]["TargetPredicateAO"];
            targetType?: string;
        };
        /**
         * @description A step in a lane executing the defined validations of a service
         * @example {
         *       "type": "service-validation",
         *       "serviceId": "1a04288d-6c85-4ba6-80d7-24ded64dd009",
         *       "parameters": {
         *         "duration": "120s"
         *       }
         *     }
         */
        ExperimentStepServiceValidationAO: {
            customLabel?: string;
            /**
             * @description Ignore any errors and failures of this single step and continue the execution of an experiment run
             * @example false
             */
            ignoreFailure?: boolean;
            /** @description Optional metric checks used to define success or failure of this step */
            metricChecks?: components["schemas"]["MetricCheckAO"][];
            /** @description Optional metric queries used of this step to filter e.g. monitoring data */
            metricQueries?: components["schemas"]["MetricQueryAO"][];
            /**
             * @description Configuration parameters of this step that are saved during experiment design and evaluated at execution time.
             * @example {
             *       "duration": "30s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            type: string;
        } & {
            /**
             * @description The name of the service to validate.
             * @example 1a04288d-6c85-4ba6-80d7-24ded64dd009
             */
            serviceId: string;
        } & {
            /**
             * @description discriminator enum property added by openapi-typescript
             * @enum {string}
             */
            type: "service-validation";
        };
        /**
         * @description A single step in a lane waiting for a specified duration.
         * @example {
         *       "type": "wait",
         *       "ignoreFailure": false,
         *       "parameters": {
         *         "duration": "10s"
         *       }
         *     }
         */
        ExperimentStepWaitAO: {
            customLabel?: string;
            /**
             * @description Ignore any errors and failures of this single step and continue the execution of an experiment run
             * @example false
             */
            ignoreFailure?: boolean;
            /** @description Optional metric checks used to define success or failure of this step */
            metricChecks?: components["schemas"]["MetricCheckAO"][];
            /** @description Optional metric queries used of this step to filter e.g. monitoring data */
            metricQueries?: components["schemas"]["MetricQueryAO"][];
            /**
             * @description Configuration parameters of this step that are saved during experiment design and evaluated at execution time.
             * @example {
             *       "duration": "30s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            type: string;
        } & {
            /**
             * @description discriminator enum property added by openapi-typescript
             * @enum {string}
             */
            type: "wait";
        };
        /**
         * @description List of experiments.
         * @example {
         *       "experiments": [
         *         {
         *           "key": "ADM-1",
         *           "name": "Shop survives unavailability of hot-deals products"
         *         },
         *         {
         *           "key": "SHOP-2",
         *           "name": "Network latency of Message Broker doesn't interfere with Online shop"
         *         }
         *       ]
         *     }
         */
        ExperimentSummariesAO: {
            /**
             * @description List of experiment summaries
             * @example [
             *       {
             *         "key": "ADM-1",
             *         "name": "Shop survives unavailability of hot-deals products"
             *       },
             *       {
             *         "key": "SHOP-2",
             *         "name": "Network latency of Message Broker doesn't interfere with Online shop"
             *       }
             *     ]
             */
            experiments?: components["schemas"]["ExperimentSummaryAO"][];
        };
        /**
         * @description Summary of a single experiment.
         * @example {
         *       "key": "ADM-1",
         *       "name": "Shop survives unavailability of hot-deals products"
         *     }
         */
        ExperimentSummaryAO: {
            /**
             * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
             * @example ADM-2
             */
            key?: string;
            /**
             * @description Name of the experiment to easily identify the experiment
             * @example Shop survives unavailability of hot-deals products
             */
            name?: string;
        };
        /**
         * @example {
         *       "id": "6bea7aec-3572-44cf-9151-c6ada57d08ca",
         *       "version": 0,
         *       "templateTitle": "HTTP Endpoint remains functional during Kubernetes Rollout Restart",
         *       "templateDescription": "Test if a given HTTP Endpoint remains funcitonal if a Kubernetes deployment is restarted.",
         *       "placeholders": [
         *         {
         *           "key": "HTTP_ENDPOINT",
         *           "name": "HTTP Endpoint",
         *           "description": "Which HTTP Endpoint should be checked during experiment execution?"
         *         },
         *         {
         *           "key": "DEPLOYMENT",
         *           "name": "Kubernetes Deployment",
         *           "description": "Which Kubernetes deployment do you want to restart?"
         *         },
         *         {
         *           "key": "CLUSTER",
         *           "name": "Kubernetes Cluster",
         *           "description": "In which Kubernetes cluster is the deployment deployed to?"
         *         },
         *         {
         *           "key": "NAMESPACE",
         *           "name": "Kubernetes Namespace",
         *           "description": "In which Kubernetes namespace is the deployment deployed to?"
         *         }
         *       ],
         *       "tags": [
         *         "Kubernetes"
         *       ],
         *       "lanes": [
         *         {
         *           "steps": [
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "60s",
         *                 "headers": [],
         *                 "method": "GET",
         *                 "successRate": "100",
         *                 "maxConcurrent": 5,
         *                 "followRedirects": false,
         *                 "readTimeout": "5s",
         *                 "connectTimeout": "5s",
         *                 "requestsPerSecond": 1,
         *                 "url": "[[HTTP_ENDPOINT]]",
         *                 "statusCode": "200-299"
         *               },
         *               "actionType": "com.steadybit.extension_http.check.periodically",
         *               "radius": {}
         *             }
         *           ]
         *         },
         *         {
         *           "steps": [
         *             {
         *               "type": "wait",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "10s"
         *               }
         *             },
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "wait": false
         *               },
         *               "actionType": "com.steadybit.extension_kubernetes.rollout-restart",
         *               "radius": {
         *                 "targetType": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *                 "predicate": {
         *                   "operator": "AND",
         *                   "predicates": [
         *                     {
         *                       "key": "k8s.cluster-name",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[CLUSTER]]"
         *                       ]
         *                     },
         *                     {
         *                       "key": "k8s.namespace",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[NAMESPACE]]"
         *                       ]
         *                     },
         *                     {
         *                       "key": "k8s.deployment",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[DEPLOYMENT]]"
         *                       ]
         *                     }
         *                   ]
         *                 },
         *                 "query": null,
         *                 "percentage": 50
         *               }
         *             },
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "10m"
         *               },
         *               "actionType": "com.steadybit.extension_kubernetes.rollout-status",
         *               "radius": {
         *                 "targetType": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *                 "predicate": {
         *                   "operator": "AND",
         *                   "predicates": [
         *                     {
         *                       "key": "k8s.cluster-name",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[CLUSTER]]"
         *                       ]
         *                     },
         *                     {
         *                       "key": "k8s.namespace",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[NAMESPACE]]"
         *                       ]
         *                     },
         *                     {
         *                       "key": "k8s.deployment",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[DEPLOYMENT]]"
         *                       ]
         *                     }
         *                   ]
         *                 },
         *                 "query": null,
         *                 "percentage": 50
         *               }
         *             }
         *           ]
         *         }
         *       ],
         *       "properties": {
         *         "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
         *       },
         *       "propertiesMetadata": [
         *         {
         *           "key": "EXAMPLE_CUSTOM_PROPERTY",
         *           "required": true,
         *           "editableInExecution": false
         *         }
         *       ],
         *       "hidden": false,
         *       "created": "2024-03-22T09:24:12.802961166Z",
         *       "createdBy": {
         *         "username": "2bd1c2d7-4051-46ad-9f05-315062edd85e",
         *         "name": "Daniel",
         *         "pictureUrl": "https://s.gravatar.com/avatar/4f27f3856530f8f2e4ec050b1d594306?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fda.png"
         *       },
         *       "edited": "2024-03-22T09:24:12.802961166Z",
         *       "editedBy": {
         *         "username": "2bd1c2d7-4051-46ad-9f05-315062edd85e",
         *         "name": "Daniel",
         *         "pictureUrl": "https://s.gravatar.com/avatar/4f27f3856530f8f2e4ec050b1d594306?s=480&r=pg&d=https%3A%2F%2Fcdn.auth0.com%2Favatars%2Fda.png"
         *       }
         *     }
         */
        ExperimentTemplateAO: {
            /**
             * Format: date-time
             * @description Timestamp when the experiment template was created
             * @example 2023-01-01T09:00:00Z
             */
            created: string;
            createdBy: components["schemas"]["UserSummaryAO"];
            /**
             * Format: date-time
             * @description Timestamp when the experiment template was edited the last time
             * @example 2023-01-01T09:00:00Z
             */
            edited: string;
            editedBy: components["schemas"]["UserSummaryAO"];
            /**
             * @description Name of the experiment created by this template. If omitted, the name needs to be added when the template is used.
             * @example Shop survives unavailability of database
             */
            experimentName?: string | null;
            /**
             * @description Should the experiment template be hidden
             * @example false
             */
            hidden?: boolean;
            /**
             * @description The hypothesis that is validated by the experiment
             * @example System is able to survive a latency in the network of 1500ms
             */
            hypothesis?: string;
            /** Format: uuid */
            id?: string | null;
            /**
             * @description The lanes (steps executed in parallel) in the experiment template. Each lane consists of multiple steps that are executed sequential per lane.
             * @example [
             *       {
             *         "steps": [
             *           {
             *             "type": "action",
             *             "ignoreFailure": false,
             *             "parameters": {
             *               "duration": "30s"
             *             },
             *             "actionType": "com.steadybit.extension_host.stress-cpu",
             *             "radius": {
             *               "targetType": "com.steadybit.extension_container.container",
             *               "predicate": {
             *                 "operator": "AND",
             *                 "predicates": [
             *                   {
             *                     "key": "k8s.deployment",
             *                     "operator": "EQUALS",
             *                     "values": [
             *                       "hot-deals"
             *                     ]
             *                   }
             *                 ]
             *               },
             *               "query": null,
             *               "percentage": 100
             *             }
             *           }
             *         ]
             *       }
             *     ]
             */
            lanes: components["schemas"]["ExperimentLaneAO"][];
            /** @description A list of placeholders used in this experiment template. */
            placeholders?: components["schemas"]["ExperimentTemplatePlaceholderAO"][];
            /**
             * @description The properties of the experiment
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * @description Metadata for properties used in this template.
             * @example [
             *       {
             *         "key": "EXAMPLE_CUSTOM_PROPERTY",
             *         "required": true,
             *         "editableInExecution": false
             *       }
             *     ]
             */
            propertiesMetadata?: components["schemas"]["PropertyMetadataAO"][];
            /** @description A list of tags for this experiment template. (Up to 5) */
            tags?: string[];
            /** @description A brief description what the template is doing. */
            templateDescription: string;
            /**
             * @description The title of the template
             * @example Shop survives unavailability of database
             */
            templateTitle: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /** @description A list of placeholders used in this experiment template. */
        ExperimentTemplatePlaceholderAO: {
            description: string;
            key: string;
            name: string;
        };
        /** @description List of template placeholder values */
        ExperimentTemplatePlaceholderValueAO: {
            /**
             * @description The key of a template placeholder
             * @example cluster-name
             */
            key: string;
            /**
             * @description The value of a template placeholder, can be a string, a number a boolean or an object with a given structure like `[{"key": "CLUSTER","value": "demo-cluster"}]`
             * @example prod-cluster-1
             */
            value: Record<string, never>;
        };
        /** @description Request to import templates from a hub. */
        ExperimentTemplatesImportAO: {
            /**
             * Format: uuid
             * @description ID of the hub to import templates from.
             */
            hubId: string;
            /** @description Optional list of Template IDs to import into the platform. If not provided, all templates from the hub will be imported. */
            templateIds?: string[] | null;
        };
        /**
         * @description List of experiment template summaries.
         * @example {
         *       "templates": [
         *         {
         *           "id": "f5990c81-6427-4144-8304-eda765a3f852",
         *           "templateTitle": "xxx"
         *         },
         *         {
         *           "id": "d7e65100-1d20-4980-be87-c351704910b8",
         *           "templateTitle": "yyy"
         *         }
         *       ]
         *     }
         */
        ExperimentTemplateSummariesAO: {
            /**
             * @description List of experiment template summaries.
             * @example [
             *       {
             *         "id": "f5990c81-6427-4144-8304-eda765a3f852",
             *         "templateTitle": "xxx"
             *       },
             *       {
             *         "key": "d7e65100-1d20-4980-be87-c351704910b8",
             *         "templateTitle": "yyy"
             *       }
             *     ]
             */
            templates?: components["schemas"]["ExperimentTemplateSummaryAO"][];
        };
        /**
         * @description Summary of a single experiment template.
         * @example {
         *       "id": "e50deab2-2636-4a5b-ad5a-6cf904ed56c4",
         *       "templateTitle": "HTTP Endpoint remains functional during Kubernetes Rollout Restart"
         *     }
         */
        ExperimentTemplateSummaryAO: {
            /**
             * @description Is the template currently hidden?
             * @example true
             */
            hidden?: boolean;
            /**
             * Format: uuid
             * @description Unique id that identifies the experiment template.
             * @example b9f4aae2-9b03-4ad3-a1a7-654774cc04eb
             */
            id?: string;
            /**
             * @description Description of the experiment template
             * @example Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.
             */
            templateDescription?: string;
            /**
             * @description Title of the experiment template
             * @example Shop survives unavailability of hot-deals products
             */
            templateTitle?: string;
        };
        /**
         * @description Request for getting and filtering pieces of advice.
         * @example {
         *       "environmentName": "Global",
         *       "query": "k8s.cluster-name=sandbox-demo and k8s.namespace=steadybit-demo",
         *       "offset": 0
         *     }
         */
        GetAdviceApiRequestAO: {
            /**
             * @description The name of the environment of which the pieces of Advice should be listed
             * @example Global
             */
            environmentName: string;
            /**
             * Format: int64
             * @description The offset to be returned in the paginated result set
             * @example 20
             */
            offset?: number | null;
            /**
             * @description An additional optional filter to search only for pieces of advice, whose target is included in the filter
             * @example k8s.cluster-name=prod-demo and k8s.namespace=steadybit-demo
             */
            query?: string | null;
        };
        GetLicenseSummaryAO: {
            /** Format: date-time */
            expires?: string | null;
            features?: components["schemas"]["LicenseFeatureSummaryAO"][] | null;
            license?: components["schemas"]["LicenseSummaryAO"];
            tenantKey?: string | null;
        };
        /** @description An informational or warning hint displayed to the user. */
        HintAO: {
            /** @description Content of the hint as markdown text. */
            content: string;
            /**
             * @description Type of the hint.
             * @example INFO
             */
            type: string;
        } | null;
        HubAO: {
            /**
             * Format: date-time
             * @description Timestamp when the hub was connected
             * @example 2023-01-01T09:00:00Z
             */
            created: string;
            createdBy: components["schemas"]["UserSummaryAO"];
            /**
             * Format: date-time
             * @description Timestamp when the hub was edited the last time
             * @example 2023-01-01T09:00:00Z
             */
            edited: string;
            editedBy: components["schemas"]["UserSummaryAO"];
            /**
             * @description Website address of the the hub
             * @example https://hub.steadybit.com/
             */
            hubLink?: string | null;
            /** @description Name of the hub */
            hubName: string;
            /** Format: uuid */
            id: string;
            /**
             * Format: date-time
             * @description Timestamp of last change as defined in the repository content
             */
            lastRepositoryChange?: string | null;
            /**
             * Format: date-time
             * @description Timestamp of last hub synchronization
             */
            lastSync?: string | null;
            /**
             * @description HTTP address of the the hub's repository
             * @example https://github.com/steadybit/reliability-hub-db
             */
            repositoryUrl: string;
            /** @description Last synchronization error description, if an error occurred */
            syncError?: string | null;
            /** @description List of templates published in the hub. */
            templates: components["schemas"]["ExperimentTemplateSummaryAO"][];
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version: number;
        };
        HubConnectionCheckAO: {
            /**
             * @description HTTP address of the the hub's repository
             * @example https://github.com/steadybit/reliability-hub-db
             */
            repositoryUrl: string;
        };
        HubConnectionCheckResponseAO: {
            error?: string | null;
        };
        /**
         * @description List of all hubs. Fetch a single hub by `id` to get more information.
         * @example {
         *       "hubs": [
         *         {
         *           "id": "1b267dc1-5f4e-4803-894d-92ecd9b83413",
         *           "hubName": "Hub Name"
         *         }
         *       ]
         *     }
         */
        HubSummariesAO: {
            hubs?: components["schemas"]["HubSummaryAO"][];
        };
        /**
         * @description Summary containing the most important hub details.
         * @example {
         *       "id": "1b267dc1-5f4e-4803-894d-92ecd9b83413",
         *       "hubName": "Hub Name"
         *     }
         */
        HubSummaryAO: {
            /** @description Name of the hub */
            hubName: string;
            /** Format: uuid */
            id: string;
        };
        /**
         * @description Request to invite users to the platform
         * @example {
         *       "email": "aa@bb.com",
         *       "role": "USER",
         *       "teamKey": "TST"
         *     }
         */
        InvitationAO: {
            /** Format: email */
            email: string;
            /** @enum {string} */
            role?: "ADMIN" | "USER";
            teamKey?: string | null;
        };
        /**
         * @description Request to invite users to the platform
         * @example {
         *       "invitations": [
         *         {
         *           "email": "aa@bb.com",
         *           "role": "ADMIN",
         *           "teamKey": "ADM"
         *         }
         *       ]
         *     }
         */
        InviteUsersRequestAO: {
            invitations: components["schemas"]["InvitationAO"][];
        };
        /**
         * @description Determines the current status of the kill switch (emergency stop). If the kill switch is active, all experiments are cancelled immediately and no new experiments can be executed.
         * @example {
         *       "active": "true",
         *       "engagedBy": "71ab0180-8abc-4d30-8acb-6aa024e3065f",
         *       "engaged": "2023-01-01T09:00:00Z"
         *     }
         */
        KillswitchAO: {
            /**
             * @description Determines whether the kill switch is currently active / engaged.
             * @example true
             */
            active?: boolean;
            /**
             * Format: date-time
             * @description Time at which the kill switch was activated / engaged
             * @example 2023-01-01T09:00:00Z
             */
            engaged?: string;
            /**
             * @description Username (internal identifier of Steadybit) of the user that has activated / engaged the kill switch
             * @example 13av2737-b318-4048-a79d-4789d645bc31
             */
            engagedBy?: string;
            engagedByDetails?: components["schemas"]["UserSummaryAO"];
        };
        /** @description A saved view of the explorer landscape. */
        LandscapeViewAO: {
            colorBy?: components["schemas"]["LandscapeViewColorByAO"];
            /**
             * @description Description of the saved view.
             * @example All shop workloads grouped by namespace.
             */
            description?: string;
            /**
             * @description Name of the environment the view is scoped to.
             * @example Global
             */
            environment?: string;
            /**
             * @description Explorer filter query narrowing the targets shown on the landscape.
             * @example k8s.namespace="shop"
             */
            filterQuery?: string;
            /** @description Ordered list of group-by dimensions the targets are grouped by, each with its own advanced configuration. */
            groupBy?: components["schemas"]["LandscapeViewGroupByAO"][];
            /**
             * Format: uuid
             * @description Unique identifier of the saved view.
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id?: string;
            /**
             * Format: date-time
             * @description Point in time the saved view was last updated.
             * @example 2026-07-23T10:15:30Z
             */
            lastUpdated?: string;
            /**
             * @description Title of the saved view.
             * @example Kubernetes by namespace
             */
            name?: string;
            /**
             * @description Whether reliability advice is shown on the landscape.
             * @example false
             */
            showAdvice?: boolean;
            /**
             * @description Attribute key the size of a target is derived from.
             * @example k8s.container.cpu.limit
             */
            sizeBy?: string;
            /**
             * @description Key of the team the saved view belongs to.
             * @example ADM
             */
            team: string;
        };
        /** @description The color-by dimension: the attribute the targets are colored by, together with its advanced configuration. */
        LandscapeViewColorByAO: {
            /**
             * @description Attribute key the color of a target is derived from.
             * @example k8s.namespace
             */
            attribute?: string;
            /** @description Buckets that map specific attribute values to named color groups. */
            mappings?: components["schemas"]["LandscapeViewMappedGroupingAO"][];
            /**
             * @description Explicit color overrides keyed by the color-by attribute value. Each color must be one of the predefined landscape colors.
             * @example {
             *       "shop": "GREEN"
             *     }
             */
            overrides?: {
                [key: string]: "RED" | "ORANGE_DARK" | "ORANGE_LIGHT" | "YELLOW" | "GREEN" | "TIFFANY" | "TEAL" | "BLUE_LIGHT" | "BLUE" | "PLUM" | "BERRIES" | "VIOLET" | "ROSE_PINK" | "PINK" | "GREY";
            } | null;
        };
        /** @description A single group-by dimension: the attribute the targets are grouped by, together with its advanced configuration. */
        LandscapeViewGroupByAO: {
            /**
             * @description Attribute key the targets are grouped by.
             * @example k8s.namespace
             */
            attribute: string;
            /** @description Buckets that map specific attribute values to named groups. */
            mappings?: components["schemas"]["LandscapeViewMappedGroupingAO"][];
            /**
             * @description Whether attribute values that are not mapped to any bucket are merged into the unknown group.
             * @example false
             */
            mergeUnmappedToUnknown?: boolean | null;
            /**
             * @description Whether an additional group collecting all targets without a value for this dimension is shown.
             * @example true
             */
            showUnknown?: boolean;
        };
        /** @description A bucket that groups multiple attribute values under a single named group. */
        LandscapeViewMappedGroupingAO: {
            /**
             * @description Attribute values that are collected into this bucket.
             * @example [
             *       "prod",
             *       "production"
             *     ]
             */
            attributeValues?: string[];
            /**
             * @description Display name of the bucket.
             * @example Production
             */
            groupName?: string;
            /**
             * Format: uuid
             * @description Unique identifier of the bucket. Optional on create/update — a new identifier is generated when omitted; supply the returned identifier to keep a bucket stable across updates.
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id?: string | null;
        };
        LicenseFeatureSummaryAO: {
            /** Format: int32 */
            hardLimit?: number | null;
            name: string;
            /** Format: int32 */
            softLimit?: number | null;
            /** @enum {string} */
            type: "SIMPLE" | "SOFT_LIMIT" | "HARD_LIMIT";
            /** Format: int32 */
            usage?: number | null;
        };
        LicenseSummaryAO: {
            /** Format: int64 */
            id?: number;
            /** @enum {string} */
            licenseType?: "NONE" | "TRIAL" | "STARTUP" | "PROFESSIONAL" | "ENTERPRISE";
            orderNumber?: string;
            /** Format: date */
            validFrom: string;
            /** Format: date */
            validTo: string;
        } | null;
        LinkCustomExperimentRequestAO: {
            /**
             * @description The category to which the experiment should be linked
             * @example Scalability
             */
            category: string;
            /**
             * @description The experiment that should be linked
             * @example ADM-18
             */
            experimentKey: string;
        };
        ListResponseCustomWebhookAO: {
            content?: components["schemas"]["CustomWebhookAO"][];
        };
        ListResponseLandscapeViewAO: {
            content?: components["schemas"]["LandscapeViewAO"][];
        };
        ListResponsePreflightActionIntegrationAO: {
            content?: components["schemas"]["PreflightActionIntegrationAO"][];
        };
        ListResponsePreflightWebhookAO: {
            content?: components["schemas"]["PreflightWebhookAO"][];
        };
        ListResponseSlackWebhookAO: {
            content?: components["schemas"]["SlackWebhookAO"][];
        };
        /**
         * @description Member of a team.
         * @example {
         *       "username": "13av2737-b318-4048-a79d-4789d645bc31",
         *       "role": "OWNER",
         *       "name": "Max Mustermann",
         *       "email": "aa@bb.com"
         *     }
         */
        MemberAO: {
            email?: string | null;
            /**
             * @description How a team or team membership is managed
             * @example MANUAL
             * @enum {string}
             */
            managedBy?: "MANUAL" | "OIDC" | "LDAP";
            /**
             * @description Name of the user
             * @example Jane Doe
             */
            name?: string;
            pictureUrl?: string | null;
            /**
             * @description Role of the team member
             * @example OWNER
             * @enum {string}
             */
            role: "MEMBER" | "OWNER";
            /**
             * @description Username of the user, internal identifier of Steadybit
             * @example 13av2737-b318-4048-a79d-4789d645bc31
             */
            username: string;
        };
        /**
         * @description Add a Member to a Team by providing the username or the email.
         * @example {
         *       "members": [
         *         {
         *           "username": "example",
         *           "email": "example@example.com",
         *           "role": "MEMBER"
         *         }
         *       ]
         *     }
         */
        MemberUpdateAO: {
            /**
             * @description E-mail of the user, unique within Steadybit
             * @example example@example.com
             */
            email?: string;
            /**
             * @description Role of the team member
             * @example OWNER
             * @enum {string}
             */
            role: "MEMBER" | "OWNER";
            /**
             * @description Username of the user, internal identifier of Steadybit
             * @example 13av2737-b318-4048-a79d-4789d645bc31
             */
            username?: string;
        };
        /** @description Optional metric checks used to define success or failure of this step */
        MetricCheckAO: {
            a: components["schemas"]["MetricValueAO"];
            b?: components["schemas"]["MetricValueAO"] | components["schemas"]["ScalarValueAO"] | components["schemas"]["VariableValueAO"];
            /** @enum {string} */
            condition: "LT" | "LTE" | "EQ" | "NEQ" | "GT" | "GTE" | "DATA_SERIES_PRESENCE";
            /** Format: uuid */
            id: string;
        };
        /** @description Optional metric queries used of this step to filter e.g. monitoring data */
        MetricQueryAO: {
            /** Format: uuid */
            id: string;
            label: string;
            parameters: {
                [key: string]: unknown;
            };
        };
        MetricValueAO: {
            type: "MetricValueAO";
        } & (Omit<components["schemas"]["ComparableValueAO"], "type"> & {
            metric?: {
                [key: string]: string;
            };
            name?: string | null;
        });
        NegationTargetPredicateAO: Record<string, never> & {
            not?: components["schemas"]["TargetPredicateAO"];
        };
        OptionAO: {
            attribute?: string;
            label?: string;
            value?: string;
        };
        PagedResponseAOAccessTokensPageItemAO: {
            items?: components["schemas"]["AccessTokensPageItemAO"][];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        PagedResponseAOAccessTokensPageItemV2AO: {
            items?: components["schemas"]["AccessTokensPageItemV2AO"][];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        PagedResponseAOExperimentExecutionPageItemAO: {
            items?: components["schemas"]["ExperimentExecutionPageItemAO"][];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        PagedResponseAOPropertyAssociationAO: {
            items?: components["schemas"]["PropertyAssociationAO"][];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        PagedResponseAOPropertyDefinitionAO: {
            items?: components["schemas"]["PropertyDefinitionAO"][];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        PagedResponseAOServiceExperimentAO: {
            items?: components["schemas"]["ServiceExperimentAO"][];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        PagedResponseAOServiceProfileAO: {
            items?: components["schemas"]["ServiceProfileAO"][];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        PagedResponseAOServiceSummaryAO: {
            items?: components["schemas"]["ServiceSummaryAO"][];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        PagedResponseAOString: {
            items?: string[];
            /**
             * Format: int32
             * @description Next page to query for next page of runs or null if there are none.
             * @example 4
             */
            nextPage?: number | null;
            /**
             * Format: int64
             * @description Total amount of runs matching your query
             * @example 241
             */
            totalItems?: number;
        };
        PageRequestAO: {
            /** Format: int32 */
            page?: number;
            /** Format: int32 */
            size?: number;
        };
        /** @description Parameters that describe how to fetch metrics for this action. */
        ParameterAO: {
            acceptedFileTypes?: string[] | null;
            advanced?: boolean;
            defaultValue?: string;
            deprecated?: boolean;
            deprecationMessage?: string;
            description?: string;
            durationUnits?: string[] | null;
            hint?: components["schemas"]["HintAO"];
            label: string;
            /** Format: int32 */
            max?: number;
            /** Format: int32 */
            min?: number;
            name: string;
            options?: components["schemas"]["OptionAO"][];
            optionsOnly?: boolean;
            /** Format: int32 */
            order?: number;
            required?: boolean;
            type: string;
        };
        /**
         * @description A partial update for an experiment schedule. Only non-null fields will be updated.
         * @example {
         *       "enabled": false
         *     }
         */
        PatchExperimentScheduleAO: {
            /**
             * @description Should the experiment run if another experiment is running?
             * @example true
             */
            allowParallel?: boolean | null;
            /**
             * @description Cron expression for the experiment schedule. If provided, startAt will be cleared.
             * @example 0 15 10 ? * *
             */
            cron?: string | null;
            /**
             * @description If `false`, the schedule is deactivated and no experiment will be executed.
             * @example false
             */
            enabled?: boolean | null;
            /**
             * Format: date-time
             * @description Start date for a single execution. If provided, cron will be cleared.
             */
            startAt?: string | null;
            /**
             * @description Optional timezone for a experiment schedule. Can only be used with `cron`.
             * @example Europe/Berlin
             */
            timezone?: string | null;
            /**
             * @description Variables that will be used when the experiment will be executed. The variables will override existing environment or experiment variables. Each value is either a constant string, an array of constant strings, or a select expression object.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal"
             *     }
             */
            variables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            } | null;
        };
        /**
         * @description A pageable list of pieces of preflight actions.
         * @example {
         *       id: "com.steadybit.extension_preflight.preflightaction.check-configuration",
         *         version: "0.1.0",
         *         description: "Check if a execution of a specific experiment in a environment is permitted.",
         *         targetAttributeIncludes: [
         *           "k8s.cluster-name",
         *           "k8s.namespace"
         *         ]
         *         }
         */
        PreflightActionAO: {
            /**
             * @description The description of the preflight action
             * @example Check if a execution of a specific experiment in a environment is permitted.
             */
            description?: string;
            /**
             * @description The unique identifier of the preflight action
             * @example com.steadybit.extension_preflight.preflightaction.check-configuration
             */
            id: string;
            /**
             * @description The name of the preflight action
             * @example Check configuration
             */
            name: string;
            /**
             * @description The list of target attributes that are included in the preflight action
             * @example [
             *       "k8s.cluster-name",
             *       "k8s.namespace"
             *     ]
             */
            targetAttributeIncludes?: string[];
            /**
             * @description The version of the preflight action
             * @example 0.1.0
             */
            version: string;
        };
        /**
         * @example {
         *       "id": "ac456d58-8fb2-4df4-86d8-ca81d7562739",
         *       "version": 1,
         *       "scope": "TEAM",
         *       "team": "ADM",
         *       "name": "Example Preflight Action Integration",
         *       "preflightActionId": "com.example.preflightaction.MyPreflightAction",
         *       "inflightInterval": "10s",
         *       "inflightTimeout": "5s"
         *     }
         */
        PreflightActionIntegrationAO: {
            /**
             * Format: uuid
             * @description The id of the webhook
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id: string;
            inflightInterval?: string;
            inflightTimeout?: string;
            /**
             * @description The name of the preflightActionIntegration
             * @example Preflight PreflightActionIntegration
             */
            name: string;
            /**
             * @description The preflight action id which is used to identify the preflight action
             * @example com.example.preflightaction.MyPreflightAction
             */
            preflightActionId: string;
            /**
             * @description The scope of the preflight action integration
             * @example TEAM
             * @enum {string}
             */
            scope: "GLOBAL" | "TEAM";
            /**
             * @description The key of the team if the scope is `TEAM`
             * @example ADM
             */
            team?: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version: number;
        };
        /**
         * @example {
         *       "id": "ac456d58-8fb2-4df4-86d8-ca81d7562739",
         *       "version": 1,
         *       "scope": "TEAM",
         *       "team": "ADM",
         *       "name": "Example Preflight Action Integration",
         *       "preflightActionId": "com.example.preflightaction.MyPreflightAction",
         *       "inflightInterval": "10s",
         *       "inflightTimeout": "5s"
         *     }
         */
        PreflightActionIntegrationUpsertAO: {
            /**
             * Format: uuid
             * @description The id of the webhook or null if a new webhook should be created.
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id?: string | null;
            inflightInterval?: string;
            inflightTimeout?: string;
            /**
             * @description The name of the preflightActionIntegration
             * @example Preflight PreflightActionIntegration
             */
            name: string;
            /**
             * @description The preflight action id which is used to identify the preflight action
             * @example com.example.preflightaction.MyPreflightAction
             */
            preflightActionId: string;
            /**
             * @description The scope of the preflight action integration
             * @example TEAM
             * @enum {string}
             */
            scope: "GLOBAL" | "TEAM";
            /**
             * @description The key of the team if the scope is `TEAM`
             * @example ADM
             */
            team?: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /**
         * @description A pageable list of pieces of perflight actions.
         * @example {
         *       "totalItems": 108,
         *       "nextOffset": 3,
         *       "items": []
         *     }
         */
        PreflightActionSummaryAO: {
            items?: components["schemas"]["PreflightActionAO"][];
            /**
             * Format: int32
             * @description Next queryable offset to query for next batch of preflight actions
             * @example 21
             */
            nextOffset?: number | null;
            /**
             * Format: int64
             * @description Total amount of preflight actions
             * @example 241
             */
            totalItems?: number;
        };
        PreflightWebhookAO: {
            /**
             * @description The events that you want to intercept. Currently only `experiment.execution.preflight` is supported.
             * @example [
             *       "experiment.execution.preflight"
             *     ]
             */
            events: string[];
            /**
             * @description Additional headers to include in the webhook request.
             * @example {
             *       "X-Custom-Header": "CustomValue",
             *       "X-Another-Header": "AnotherValue"
             *     }
             */
            headers?: {
                [key: string]: string;
            };
            /**
             * Format: uuid
             * @description The id of the webhook
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id: string;
            /**
             * @description The name of the webhook
             * @example Preflight Webhook
             */
            name: string;
            /**
             * @description The scope of the webhook / integration
             * @example TEAM
             * @enum {string}
             */
            scope: "GLOBAL" | "TEAM";
            /**
             * @description If a secret is provided a signature of the body is computed using `HMAC SHA-256` and sent as `X-SB-Signature` http header. You can use this header to verify the message.
             * @example secret123!!
             */
            secret?: string;
            /**
             * @description The body size can get very large as we include all target attributes for each target of your experiments. When having experiments with many targets, it might be useful to filter the attributes to only include the ones you are interested in. You can use the wildcard character '*' to match all attributes or a comma-separated-list of attribute-names. If the field is empty, no attributes will be included.
             * @example [
             *       "k8s.cluster-name",
             *       "k8s.deployment"
             *     ]
             */
            targetAttributeIncludes: string[];
            /**
             * @description The key of the team if the scope is `TEAM`
             * @example ADM
             */
            team?: string;
            /**
             * @description The URL of the webhook
             * @example https://example.com/webhook
             */
            url: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version: number;
        };
        /**
         * @description Optional preflight webhook response body that can be used to deliver a message of a preflight webhook. This message is used as failure / error / success reason depending on the HTTP API preflight response status code.
         * @example {
         *       "message": "Experiment executions are not allowed for this environment"
         *       "modifications": [
         *         {
         *           "type": "set_property_value",
         *           "propertyKey": "approvedBy",
         *           "value": "Daniel"
         *         }
         *       ]
         *     }
         */
        PreflightWebhookResponseAO: {
            /**
             * @description The message that should be used as a reason for the successful / errored / failed preflight webhook
             * @example Experiment executions are not allowed for this environment.
             */
            message?: string;
            /**
             * @description Modifications that should be applied to the execution if the preflight check was successful.
             * @example [
             *       {
             *         "type": "set_property_value",
             *         "propertyKey": "approvedBy",
             *         "value": "Daniel"
             *       },
             *       {
             *         "type": "add_value_to_list_property",
             *         "propertyKey": "observations",
             *         "value": "This looks interesting!"
             *       }
             *     ]
             */
            modifications?: (components["schemas"]["AddValueToListProperty"] | components["schemas"]["SetPropertyValue"])[];
        };
        PreflightWebhookUpsertAO: {
            /**
             * @description The events that you want to intercept. Currently only `experiment.execution.preflight` is supported.
             * @example [
             *       "experiment.execution.preflight"
             *     ]
             */
            events: string[];
            /**
             * @description Additional headers to include in the webhook request.
             * @example {
             *       "X-Custom-Header": "CustomValue",
             *       "X-Another-Header": "AnotherValue"
             *     }
             */
            headers?: {
                [key: string]: string;
            };
            /**
             * Format: uuid
             * @description The id of the webhook or null if a new webhook should be created.
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id?: string | null;
            /**
             * @description The name of the webhook
             * @example Preflight Webhook
             */
            name: string;
            /**
             * @description The scope of the webhook / integration
             * @example TEAM
             * @enum {string}
             */
            scope: "GLOBAL" | "TEAM";
            /**
             * @description If a secret is provided a signature of the body is computed using `HMAC SHA-256` and sent as `X-SB-Signature` http header. You can use this header to verify the message.
             * @example secret123!!
             */
            secret?: string;
            /**
             * @description The body size can get very large as we include all target attributes for each target of your experiments. When having experiments with many targets, it might be useful to filter the attributes to only include the ones you are interested in. You can use the wildcard character '*' to match all attributes or a comma-separated-list of attribute-names. If the field is empty, no attributes will be included.
             * @example [
             *       "k8s.cluster-name",
             *       "k8s.deployment"
             *     ]
             */
            targetAttributeIncludes: string[];
            /**
             * @description The key of the team if the scope is `TEAM`
             * @example ADM
             */
            team?: string;
            /**
             * @description The URL of the webhook
             * @example https://example.com/webhook
             */
            url: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /**
         * @description The principal that has performed the logged event
         * @example {
         *       "name": "Jane Doe",
         *       "role": "ADMIN",
         *       "email": "example@example.com",
         *       "username": "1ava2afg-xju33-4c6a-9451-2854584c15be",
         *       "principalType": "USER"
         *     }
         */
        PrincipalAL: ({
            /**
             * @description The principal type that has performed the logged event
             * @example USER
             * @enum {string}
             */
            principalType: "USER" | "ACCESS_TOKEN" | "BATCH_JOB";
        } & (components["schemas"]["AccessTokenPrincipalAL"] | components["schemas"]["BatchPrincipalAL"] | components["schemas"]["UserPrincipalAL"])) | null;
        /**
         * @description A property association.
         * @example {
         *       "id": "2v1av42-e525-4c00-a13a-1ac32d170724",
         *       "key": "RESULT_COLOR",
         *       "editableInExecution": true,
         *       "required": true,
         *       "version": 1
         *     }
         */
        PropertyAssociationAO: {
            /**
             * @description Always defined to either `EXPERIMENT` for experiment design or run related associations or `SERVICE` for service-associations. Only for the former, an `experimentKey` can be defined and only for the latter, a `serviceId` can be defined
             * @default EXPERIMENT
             * @example EXPERIMENT
             * @enum {string}
             */
            associationType: "EXPERIMENT" | "SERVICE";
            /**
             * @description Is the property editable in the execution view. Only used when `associationType` is set to `EXPERIMENT`.
             * @example true
             */
            editableInExecution?: boolean;
            /**
             * @description The key of the associated experiment. When `associationType` is set to `EXPERIMENT` and `experimentKey` is `null`, it is associated to ALL experiment designs. Can't be changed during updates.
             * @example EXP-1
             */
            experimentKey?: string | null;
            /**
             * Format: uuid
             * @description Id of the Property-Association.
             */
            id: string;
            /**
             * @description The key of the property definition
             * @example RESULT_COLOR
             */
            key: string;
            /**
             * @description Is the value required?
             * @example true
             */
            required?: boolean;
            /**
             * Format: uuid
             * @description The serviceId of the associated service. When `associationType` is set to `SERVICE` and `serviceId` is `null`, it is associated to ALL services. Can't be changed during updates.
             * @example 3308b47d-5c1f-4f08-a25b-a18fc10f8a56
             */
            serviceId?: string | null;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version: number;
        };
        /**
         * @description Definition of a property definition that can be associated.
         * @example {
         *       "key": "RESULT_COLOR",
         *       "label": "Result Color",
         *       "description": "How would you describe the result of your experiment, thinking in beautiful colors?",
         *       "dataType": "ENUM",
         *       "enumValues": [
         *         "RED",
         *         "GREEN",
         *         "BLUE"
         *       ],
         *       "version": 1
         *     }
         */
        PropertyDefinitionAO: {
            /**
             * @description The data type of the property
             * @example STRING
             * @enum {string}
             */
            dataType: "STRING" | "STRING_LIST" | "ENUM" | "ENUM_LIST" | "NUMBER" | "NUMBER_LIST" | "MARKDOWN" | "BOOLEAN" | "DATE" | "LINK" | "LINK_LIST";
            /**
             * @description The text describing the property.
             * @example How would you describe the result of your experiment, thinking in beautiful colors?
             */
            description?: string | null;
            /**
             * @description Valid values if the dataType `ENUM` is used
             * @example [
             *       "RED",
             *       "GREEN",
             *       "BLUE"
             *     ]
             */
            enumValues?: string[];
            /**
             * @description The unique key of the property definition
             * @example RESULT_COLOR
             */
            key: string;
            /**
             * @description The label shown in the ui for this property
             * @example Result color
             */
            label: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version: number;
        };
        /**
         * @description Insert or update the experiment template in Steadybit. The `id` will be used to identify whether the template already exists and should be updated or newly inserted.
         * @example {
         *       "key": "EXAMPLE_CUSTOM_PROPERTY",
         *       "required": true,
         *       "editableInExecution": false
         *     }
         */
        PropertyMetadataAO: {
            editableInExecution?: boolean;
            key: string;
            required?: boolean;
        };
        QueryLanguagePredicateAO: Record<string, never> & {
            query: string;
        };
        RecreateAccessTokenRequestV2AO: {
            /**
             * Format: date-time
             * @description New expiration date for the recreated token.
             * @example 2027-01-01T00:00:00Z
             */
            expiresAt?: string;
        };
        /** @description Filter for time-series report data. */
        ReportFilterAO: {
            /**
             * Format: date
             * @description Start date of the report range (inclusive).
             * @example 2026-01-01
             */
            from: string;
            /**
             * @description The time bucket granularity for report aggregation.
             * @example MONTHLY
             * @enum {string}
             */
            rollup?: "MONTHLY" | "DAILY";
            /**
             * Format: date
             * @description End date of the report range (inclusive).
             * @example 2026-03-01
             */
            to: string;
        };
        ScalarValueAO: {
            type: "ScalarValueAO";
        } & (Omit<components["schemas"]["ComparableValueAO"], "type"> & {
            /** Format: double */
            value?: number;
        });
        SelectExpressionAO: {
            attribute: string;
            /** Format: int32 */
            count?: number | null;
            filter?: string;
            /** @enum {string} */
            mode: "fixed" | "percent";
            /** Format: int32 */
            percent?: number | null;
            /**
             * @description Evaluation scope of a dynamic value, **required for service variables** and rejected for all other variables (environment, experiment, schedule, execution overrides). `service` samples from the service's own targets (its environment narrowed by the service's target query); `environment` samples from the whole environment the service lives in. There is no default — a service variable must state its scope explicitly.
             * @example service
             * @enum {string|null}
             */
            scope?: "service" | "environment" | "service" | "environment" | null;
            targetType: string;
            type?: string;
        };
        /**
         * @example {
         *       "id": "2v1av42-e525-4c00-a13a-1ac32d170724",
         *       "version": 1,
         *       "name": "shopping-service",
         *       "environment": "Global",
         *       "team": "ADM",
         *       "logoId": "service-router",
         *       "logoColor": "blue",
         *       "query": "aws.account=\"123\" OR aws.account=\"456\"",
         *       "validations": [
         *         {
         *           "type": "action",
         *           "parameters": {
         *             "url": "https://my-service/health",
         *             "method": "GET"
         *           },
         *           "actionType": "com.steadybit.extension_http.check.periodically"
         *         }
         *       ],
         *       "serviceProfile": "Steadybit Starter",
         *       "variables": {
         *         "httpEndpoint": "http://prod.shop.products.internal",
         *         "targets": {
         *           "type": "select",
         *           "targetType": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *           "attribute": "k8s.deployment",
         *           "filter": "k8s.namespace=\"shop\"",
         *           "mode": "fixed",
         *           "count": 1
         *         }
         *       },
         *       "created": "2023-01-01T09:00:00Z",
         *       "createdBy": {
         *         "username": "ag1hb7ap-d299-47ab-998f-c2a53b433820",
         *         "name": "Manuel",
         *         "pictureUrl": "https://.../picture.png"
         *       },
         *       "edited": "2023-01-01T09:00:00Z",
         *       "editedBy": {
         *         "username": "ag1hb7ap-d299-47ab-998f-c2a53b433820",
         *         "name": "Manuel",
         *         "pictureUrl": "https://.../picture.png"
         *       }
         *     }
         */
        ServiceAO: {
            /**
             * Format: date-time
             * @description Timestamp when the service was created
             * @example 2023-01-01T09:00:00Z
             */
            created: string;
            createdBy: components["schemas"]["UserSummaryAO"];
            /**
             * Format: date-time
             * @description Timestamp when the service was edited the last time
             * @example 2023-01-01T09:00:00Z
             */
            edited: string;
            editedBy: components["schemas"]["UserSummaryAO"];
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment: string;
            /**
             * Format: uuid
             * @description The unique id of the service
             */
            id?: string;
            /**
             * @description Color scheme of the logo used to identify the service in the Platform UI
             * @default blue
             * @example orangeLight
             */
            logoColor: string;
            /**
             * @description Identifier of the logo used to identify the service in the Platform UI
             * @default service
             * @example service-router
             */
            logoId: string;
            /**
             * @description The name of the service
             * @example calculator-service
             */
            name: string;
            /**
             * @description The properties of the service
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "I like this service!"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * @description Query-Language predicate, specifies the targets belonging to this Service
             * @example aws.account="123" OR aws.account="456"
             */
            query: string;
            /**
             * @description Name of the service profile that should be used for this service
             * @example Steadybit provided
             */
            serviceProfile: string;
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            team: string;
            /** @description List of validations to be executed against the service */
            validations: components["schemas"]["ExperimentStepActionAO"][];
            /**
             * @description Variables owned by the service. Each value is either a constant string, an array of constant strings, or a select expression object. A select-expression value **must set `scope`** (`service` or `environment`) — the request is rejected otherwise. On `POST /api/services` (upsert): omitting this field leaves existing variables untouched, an empty object removes all of them.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal",
             *       "targetServices": [
             *         "gateway",
             *         "hot-deals",
             *         "fashion-bestseller"
             *       ]
             *     }
             */
            variables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * Format: int32
             * @description Version for optimistic locking
             * @example 1
             */
            version?: number;
        };
        ServiceExperimentAO: {
            /**
             * @description The type of the association
             * @example PROVIDED
             * @enum {string}
             */
            associationType: "PROVIDED" | "CUSTOM";
            /**
             * @description The category of the experiment.
             * @example Scalability
             */
            category: string;
            /**
             * @description The key of the experiment. If type is PROVIDED and the experiment has not been created, the experimentKey will be null
             * @example EX-754
             */
            experimentKey?: string | null;
            /**
             * Format: uuid
             * @description The id of the experiment template if type is PROVIDED and the experiment has not been created.
             * @example 6bea7aec-3572-44cf-9151-c6ada57d08ca
             */
            templateId?: string | null;
        };
        /** @description A service profile that groups experiment templates by category */
        ServiceProfileAO: {
            /**
             * Format: date-time
             * @description Timestamp when the profile was created
             * @example 2023-01-01T09:00:00Z
             */
            created: string;
            /**
             * @description Username of the user that created the profile
             * @example admin@example.com
             */
            createdBy: string;
            /**
             * @description Whether this is the default profile.
             * @example true
             */
            defaultProfile: boolean;
            /**
             * Format: date-time
             * @description Timestamp when the profile was last edited
             * @example 2023-01-01T09:00:00Z
             */
            edited: string;
            /**
             * @description Username of the user that last edited the profile
             * @example admin@example.com
             */
            editedBy: string;
            /**
             * Format: uuid
             * @description The unique id of the profile
             */
            id: string;
            /**
             * @description The name of the profile
             * @example Default Resilience Tests
             */
            name: string;
            /**
             * @description Origin of a service profile
             * @example CUSTOM
             * @enum {string}
             */
            origin: "PROVIDED" | "CUSTOM";
            /** @description Template entries in this profile */
            templates: components["schemas"]["ServiceProfileCategoryAO"][];
            /**
             * Format: int32
             * @description Version for optimistic locking
             * @example 1
             */
            version: number;
        };
        /** @description Template entries in this profile */
        ServiceProfileCategoryAO: {
            /**
             * @description The category name
             * @example Scalability
             */
            category?: string | null;
            /**
             * @description The template IDs in this category
             * @example [
             *       "6bea7aec-3572-44cf-9151-c6ada57d08ca",
             *       "6bea7aec-3572-44cf-9151-c6ada57d08cb"
             *     ]
             */
            templateIds?: string[];
        };
        /**
         * @description The risk for a given service.
         * @example {
         *       "risk": 78,
         *       "categoryRisks": {
         *         "Scalability": {
         *           "total": 92,
         *           "experiment": 100,
         *           "advice": 58
         *         },
         *         "Redundancy": {
         *           "total": 91,
         *           "experiment": 100,
         *           "advice": 51
         *         },
         *         "Dependency": {
         *           "total": 50,
         *           "experiment": 50,
         *           "advice": 50
         *         }
         *       },
         *       "experimentRisks": [
         *         {
         *           "experimentKey": "ADM-8",
         *           "risk": 100
         *         },
         *         {
         *           "experimentKey": "ADM-16",
         *           "risk": 100
         *         }
         *       ],
         *       "lastCalculated": "2026-03-31T10:13:38.902372Z"
         *     }
         */
        ServiceRiskAO: {
            /** @description Risk per category */
            categoryRisks?: {
                [key: string]: components["schemas"]["CategoryRiskAO"];
            };
            /** @description Risk per experiment */
            experimentRisks?: components["schemas"]["ExperimentRiskAO"][];
            /**
             * Format: date-time
             * @description Timestamp of the last risk calculation
             */
            lastCalculated?: string | null;
            /**
             * Format: int32
             * @description The overall risk for the service
             */
            risk?: number;
        };
        /** @description Filter for service risk report data, optionally scoped to specific teams, environments, and services. */
        ServiceRiskReportFilterAO: {
            /** @description Restrict results to services whose categoryRisks map contains any of the given category keys. */
            categoryKeys?: string[] | null;
            /** @description Restrict results to the given environment IDs. If not provided, all environments are included. */
            environmentIds?: string[] | null;
            /**
             * Format: date
             * @description Start date of the report range (inclusive).
             * @example 2026-01-01
             */
            from: string;
            /**
             * @description The time bucket granularity for report aggregation.
             * @example MONTHLY
             * @enum {string}
             */
            rollup?: "MONTHLY" | "DAILY";
            /** @description Restrict results to the given service IDs. If not provided, all services are included. */
            serviceIds?: string[] | null;
            /** @description Filter on the service's enum/enum-list custom property values. Map of property key and values; values are OR within a key, AND across keys. */
            serviceProperties?: {
                [key: string]: string[];
            } | null;
            /** @description Restrict results to the given team IDs. If not provided, all teams are included. */
            teamIds?: string[] | null;
            /**
             * Format: date
             * @description End date of the report range (inclusive).
             * @example 2026-03-01
             */
            to: string;
        };
        ServiceSummaryAO: {
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment: string;
            /**
             * Format: uuid
             * @description The unique id of the service
             */
            id?: string;
            /**
             * @description Color scheme of the logo used to identify the service in the Platform UI
             * @example blue
             */
            logoColor?: string;
            /**
             * @description Identifier of the logo used to identify the service in the Platform UI
             * @example 1
             */
            logoId?: string;
            /**
             * @description The name of the service
             * @example calculator-service
             */
            name: string;
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            team: string;
        };
        /**
         * @description Set a value of an execution property
         * @example {
         *       "type": "set_property_value",
         *       "propertyKey": "approvedBy",
         *       "value": "Daniel"
         *     }
         */
        SetPropertyValue: {
            type: "SetPropertyValue";
        } & (Omit<components["schemas"]["ExecutionModification"], "type"> & {
            /**
             * @description The key of the property.
             * @example approvedBy
             */
            propertyKey: string;
            /**
             * @description The value to be set. Datatype depends on the property definition. Could be a string, a number or a list
             * @example Daniel
             */
            value: Record<string, never>;
        });
        SlackWebhookAO: {
            /**
             * @description The name of the slack channel
             * @example #steadybit-notifications
             */
            channel: string;
            /**
             * @description The events that are being sent or a list containing a single `*` if all supported event types should be used.
             *
             *     Supported Events:
             *      - "experiment.execution.requested"
             *      - "experiment.execution.created"
             *      - "experiment.execution.preflight"
             *      - "experiment.execution.completed"
             *      - "experiment.execution.failed"
             *      - "experiment.execution.errored"
             *      - "experiment.execution.canceled"
             *      - "experiment.execution.step-started"
             *      - "experiment.execution.step-completed"
             *      - "experiment.execution.step-failed"
             *      - "experiment.execution.step-errored"
             *      - "experiment.execution.step-canceled"
             *      - "experiment.execution.step-skipped"
             *      - "killswitch.engaged"
             *      - "killswitch.disengaged"
             * @example [
             *       "experiment.execution.created",
             *       "experiment.execution.completed"
             *     ]
             */
            events?: string[];
            /**
             * @description The icon URL of the slack channel, defaults to the Steadybit logo.
             * @example https://platform.steadybit.com/assets/logo512.png
             */
            iconUrl?: string;
            /**
             * Format: uuid
             * @description The id of the webhook
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id: string;
            /**
             * @description The name of the integration
             * @example Slack ACME corporation.
             */
            name: string;
            /**
             * @description The scope of the webhook / integration
             * @example TEAM
             * @enum {string}
             */
            scope: "GLOBAL" | "TEAM";
            /**
             * @description The key of the team if the scope is `TEAM`
             * @example ADM
             */
            team?: string;
            /**
             * @description The Slack webhook url
             * @example https://hooks.slack.com/services/<redacted>
             */
            url: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version: number;
        };
        SlackWebhookUpsertAO: {
            /**
             * @description The name of the slack channel
             * @example #steadybit-notifications
             */
            channel: string;
            /**
             * @description The events that are being sent or a list containing a single `*` if all supported event types should be used.
             *
             *     Supported Events:
             *      - "experiment.execution.requested"
             *      - "experiment.execution.created"
             *      - "experiment.execution.preflight"
             *      - "experiment.execution.completed"
             *      - "experiment.execution.failed"
             *      - "experiment.execution.errored"
             *      - "experiment.execution.canceled"
             *      - "experiment.execution.step-started"
             *      - "experiment.execution.step-completed"
             *      - "experiment.execution.step-failed"
             *      - "experiment.execution.step-errored"
             *      - "experiment.execution.step-canceled"
             *      - "experiment.execution.step-skipped"
             *      - "killswitch.engaged"
             *      - "killswitch.disengaged"
             * @example [
             *       "experiment.execution.created",
             *       "experiment.execution.completed"
             *     ]
             */
            events?: string[];
            /**
             * @description The icon URL of the slack channel, defaults to the Steadybit logo.
             * @example https://platform.steadybit.com/assets/logo512.png
             */
            iconUrl?: string;
            /**
             * Format: uuid
             * @description The id of the webhook or null if a new webhook should be created.
             * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
             */
            id?: string | null;
            /**
             * @description The name of the integration
             * @example Slack ACME corporation.
             */
            name: string;
            /**
             * @description The scope of the webhook / integration
             * @example TEAM
             * @enum {string}
             */
            scope: "GLOBAL" | "TEAM";
            /**
             * @description The key of the team if the scope is `TEAM`
             * @example ADM
             */
            team?: string;
            /**
             * @description The Slack webhook url
             * @example https://hooks.slack.com/services/<redacted>
             */
            url: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /**
         * @description Summary of a single advice for the referenced target
         * @example {
         *       "type": "com.steadybit.extension_kubernetes.advice.k8s-cpu-limit",
         *       "label": "Limit CPU Resources",
         *       "tags": [
         *         "kubernetes",
         *         "limit",
         *         "cpu"
         *       ],
         *       "status": "Validation needed",
         *       "summary": "You already took action and configured a CPU limit. Validate your configuration via an experiment."
         *     }
         */
        TargetAdviceAdvicePartAO: {
            /**
             * @description Human readable label of the advice
             * @example Limit CPU Resources
             */
            label?: string;
            /**
             * @description Current status of the advice applied to the referenced target. One of 'Action Needed', 'Validation needed', 'Implemented'.
             * @example Validation needed
             */
            status?: string;
            /**
             * @description Summary of the advice to describe the current status and next step.
             * @example You already took action and configured a CPU limit. Validate your configuration via an experiment
             */
            summary: string;
            /**
             * @description Tags associated to the advice definition
             * @example [
             *       "AWS",
             *       "Kubernetes"
             *     ]
             */
            tags?: string[];
            /**
             * @description Identifier of the advice definition that is applied to the target
             * @example com.steadybit.extension_kubernetes.advice.k8s-cpu-limit
             */
            type?: string;
        };
        /**
         * @description A pageable list of pieces of advice.
         * @example {
         *       "target": {
         *         "type": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *         "reference": "prod-demo/steadybit-demo/gateway",
         *         "label": "gateway"
         *       },
         *       "advice": {
         *         "type": "com.steadybit.extension_kubernetes.advice.k8s-cpu-limit",
         *         "label": "Limit CPU Resources",
         *         "tags": [
         *           "kubernetes",
         *           "limit",
         *           "cpu"
         *         ],
         *         "status": "Validation needed",
         *         "summary": "You already took action and configured a CPU limit. Validate your configuration via an experiment."
         *       },
         *       "url": "https://platform.steadybit.com/permalink/advice/eyAiZW52..."
         *     }
         */
        TargetAdviceAO: {
            advice?: components["schemas"]["TargetAdviceAdvicePartAO"];
            target?: components["schemas"]["TargetAdviceTargetPartAO"];
            /**
             * @description URL to see all details to this advice for this target
             * @example https://platform.steadybit.com/permalink/advice/eyAiZW52...
             */
            url?: string | null;
        };
        /**
         * @description A reference to identify the target for a given advice
         * @example {
         *       "type": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *       "reference": "prod-demo/steadybit-demo/gateway",
         *       "label": "gateway"
         *     }
         */
        TargetAdviceTargetPartAO: {
            /**
             * @description Human readable identifier to be displayed for the target
             * @example gateway
             */
            label: string;
            /**
             * @description Unique stable identifier of the target for this target type
             * @example prod-demo/steadybit-demo/gateway
             */
            reference: string;
            /**
             * @description Target type of the referenced target
             * @example com.steadybit.extension_kubernetes.kubernetes-deployment
             */
            type: string;
        };
        TargetAgentIdPredicateAO: Record<string, never> & {
            /** Format: uuid */
            agentId: string;
        };
        TargetAO: {
            /**
             * @description The ID of the agent this target belongs to
             * @example 019504aa-0f60-781b-862d-9763010d5948
             */
            agentId?: string;
            /**
             * @description The attributes for this target. A key may be associated multiple time to a single target.
             * @example [
             *       {
             *         "key": "container.port",
             *         "value": "51152:2376"
             *       },
             *       {
             *         "key": "container.engine",
             *         "value": "docker"
             *       }
             *     ]
             */
            attributes?: components["schemas"]["AttributeAO"][];
            /**
             * @description The name of the target
             * @example fashion-bestseller
             */
            name?: string;
            /**
             * @description The name of the target
             * @example com.steadybit.extension_kubernetes.kubernetes-deployment
             */
            type?: string;
            /**
             * Format: int64
             * @description The version of the target, will be increased by every update via the agent.
             * @example 0
             */
            version?: number;
        };
        TargetAttributeKeyCountPredicateAO: Record<string, never> & {
            key: string;
            value: string;
            /** @enum {string} */
            valueCountOperator: "EQUAL" | "NOT_EQUAL" | "GREATER_THAN" | "GREATER_THAN_OR_EQUAL" | "LESS_THAN" | "LESS_THAN_OR_EQUAL";
        };
        TargetAttributeKeyPredicateAO: Record<string, never> & {
            key: string;
            operator: string;
        };
        TargetAttributeKeyPresencePredicateAO: Record<string, never> & {
            key: string;
            /** @enum {string} */
            presenceOperator: "PRESENT" | "NOT_PRESENT";
        };
        TargetAttributeKeyValuePredicateAO: Record<string, never> & {
            key: string;
            operator: string;
            values: string[];
        };
        /**
         * @description A target that is expected to be effected by this action.
         * @example {
         *       "type": "com.steadybit.extension_container.container",
         *       "name": "docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea",
         *       "state": "COMPLETED",
         *       "attributes": [
         *         {
         *           "key": "container.port",
         *           "value": "51152:2376"
         *         },
         *         {
         *           "key": "container.engine",
         *           "value": "docker"
         *         },
         *         {
         *           "key": "container.host/name",
         *           "value": "docker-desktop/minikube"
         *         },
         *         {
         *           "key": "container.host",
         *           "value": "docker-desktop"
         *         }
         *       ]
         *     }
         */
        TargetExecutionAO: {
            /**
             * @description The agent that processed this target-action command and forwarded it to the proper extension instance.
             * @example prod-demo/steadybit-agent/steadybit-agent-0
             */
            agentHostname?: string;
            /**
             * @description List of artifact identifiers that are associated with this target execution
             * @example [
             *       "jmeter-report.zip",
             *       "system-metrics.csv"
             *     ]
             */
            artifacts?: string[];
            /**
             * @description A set of attributes that have been discovered for this target. A key may be associated multiple time to a single target.
             * @example [
             *       {
             *         "key": "container.port",
             *         "value": "51152:2376"
             *       },
             *       {
             *         "key": "container.engine",
             *         "value": "docker"
             *       }
             *     ]
             */
            attributes?: components["schemas"]["AttributeAO"][];
            /**
             * Format: uuid
             * @description Unique identifier of this target execution
             * @example 019aba52-558d-7d44-b793-92839f3c3152
             */
            id?: string;
            /**
             * @description Identifier of the target that is expected to be effected
             * @example docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea
             */
            name?: string;
            /**
             * @description Reason on a per target-level why the experiment failed or errored. If this step didn't failed or errored (`state != 'FAILED' and state != 'ERRORED') the reason is `null`.
             * @example Failed to start Stop Container (com.steadybit.extension_container.stop)
             */
            reason?: string;
            /**
             * @description Optional additional reason details on a per target-level why the the experiment failed or errored.
             * @example Could not read state of target container: exit status 1 (time="2023-09-29T12:41:32Z" level=error msg="container does not exist"
             */
            reasonDetails?: string;
            /**
             * @description The source (i.e. call to the extension) that caused the step to error or fail.
             * @example POST http://11.20.86.255:9093/com.steadybit.extension_container.container_stop/prepare
             */
            source?: string;
            /**
             * @description State of this specific step on a per target-level.
             * @example COMPLETED
             */
            state?: string;
            summary?: components["schemas"]["TargetExecutionSummaryAO"];
            /**
             * @description Type of the target that is expected to be effected
             * @example container
             */
            type?: string;
        };
        /**
         * @description An attributes (key-value-pair) that is associated to a target
         * @example {
         *       "key": "container.port",
         *       "value": "51152:2376"
         *     }
         */
        TargetExecutionAOAttributeAO: {
            /**
             * @description The key of the attribute, may be associated multiple times to the same target
             * @example container.engine
             */
            key: string;
            /**
             * @description The value of the attribute
             * @example docker
             */
            value: string;
        };
        /**
         * @description A summary for a target execution
         * @example {
         *       "text": "Hello world!",
         *       "level": "INFO"
         *     }
         */
        TargetExecutionSummaryAO: {
            level?: string;
            text?: string;
        };
        TargetNamePredicateAO: Record<string, never> & {
            name: string;
        };
        /**
         * @description Query defining the overall superset of targets being effected
         * @example [
         *       {
         *         "key": "container.host/name",
         *         "operator": "EQUALS",
         *         "values": [
         *           "docker-desktop/minikube"
         *         ]
         *       }
         *     ]
         */
        TargetPredicateAO: components["schemas"]["NegationTargetPredicateAO"] | components["schemas"]["QueryLanguagePredicateAO"] | components["schemas"]["TargetAgentIdPredicateAO"] | components["schemas"]["TargetAttributeKeyCountPredicateAO"] | components["schemas"]["TargetAttributeKeyPredicateAO"] | components["schemas"]["TargetAttributeKeyPresencePredicateAO"] | components["schemas"]["TargetAttributeKeyValuePredicateAO"] | components["schemas"]["TargetNamePredicateAO"] | components["schemas"]["TargetTypePredicateAO"];
        /** @description A predefined target predicate template that users can apply when configuring an action. */
        TargetPredicateTemplateAO: {
            /** @description Description of the template. */
            description?: string;
            /** @description Display name of the template. */
            name?: string;
            /** @description Query language template. */
            template: string;
        };
        TargetSelectorAO: {
            targetQuery?: string;
            type: string;
        };
        TargetStatsRequest: {
            predicate?: components["schemas"]["TargetPredicateAO"];
            /**
             * @description Alternative to `predicate`. If both `query` and `predicate` will be provided, `query` will override the `predicate`.
             * @example (aws.account="123" OR aws.account="456"
             */
            query?: string | null;
        };
        TargetTypePredicateAO: Record<string, never> & {
            types: string[];
        };
        /**
         * @description The team in which the event was triggered
         * @example {
         *       "id": "a2167b29-e73b-4445-8468-4670a0b459b3",
         *       "key": "ADMIN",
         *       "name": "Administrators"
         *     }
         */
        TeamAL: {
            /** Format: uuid */
            id: string;
            key: string;
            name: string;
        } | null;
        /**
         * @description A team that is uniquely identified via it's teamKey and has members, allowed environments and actions.
         * @example {
         *       "id": "71ab0180-8abc-4d30-8acb-6aa024e3065f",
         *       "key": "ADM",
         *       "name": "Administrators",
         *       "version": 1,
         *       "logoId": "1",
         *       "logoColor": "cyanDark",
         *       "allowedActions": [
         *         "com.steadybit.extension_host.host.stress-cpu"
         *       ],
         *       "allowedEnvironments": [
         *         "Global"
         *       ],
         *       "members": [
         *         {
         *           "username": "auth0|11a9315afc84590069cd53b2",
         *           "role": "OWNER"
         *         },
         *         {
         *           "username": "auth0|13b1s51vg184590069cd51ab",
         *           "role": "MEMBER"
         *         },
         *         {
         *           "username": "auth0|1va2g15afc84590069cd53c3",
         *           "role": "MEMBER"
         *         }
         *       ]
         *     }
         */
        TeamAO: {
            /**
             * @description Set of allowed actions that can be used in an experiment of this team
             * @example [
             *       "com.steadybit.extension_host.host.stress-cpu"
             *     ]
             */
            allowedActions: string[];
            /**
             * @description Set of allowed environments, identified via name
             * @example [
             *       "Global"
             *     ]
             */
            allowedEnvironments: string[];
            /** @description An optional description of a team */
            description?: string;
            /** Format: uuid */
            id?: string;
            /**
             * @description Unique identifier of a team
             * @example ADM
             */
            key: string;
            /**
             * @description Color scheme of the logo used to identify the team in the Platform UI
             * @example cyanDark
             */
            logoColor?: string;
            /**
             * @description Identifier of the logo used to identify the team in the Platform UI
             * @example 1
             */
            logoId?: string;
            /**
             * @description How a team or team membership is managed
             * @example MANUAL
             * @enum {string}
             */
            managedBy?: "MANUAL" | "OIDC" | "LDAP";
            /**
             * @description Members that are associated to this team
             * @example {
             *       "username": "auth0|11a9315afc84590069cd53b2",
             *       "role": "OWNER"
             *     }
             */
            members: components["schemas"]["MemberAO"][];
            /**
             * @description Name of a team
             * @example ADMIN
             */
            name: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /**
         * @description Environment assigned to a team.
         * @example {
         *       "name": "Global"
         *     }
         */
        TeamEnvironmentAO: {
            /**
             * @description Name of the environment.
             * @example Global
             */
            name: string;
        };
        /**
         * @description List of environments that are assigned to this team.
         * @example {
         *       "environments": [
         *         {
         *           "name": "Global"
         *         },
         *         {
         *           "name": "Shop Production"
         *         }
         *       ]
         *     }
         */
        TeamEnvironmentsAO: {
            /**
             * @description Environments that are assigned to this team
             * @example [
             *       {
             *         "name": "Global"
             *       },
             *       {
             *         "name": "Shop Production"
             *       }
             *     ]
             */
            environments: components["schemas"]["TeamEnvironmentAO"][];
        };
        /**
         * @description Update request to change the environments of a specific team.
         * @example {
         *       "environments": [
         *         {
         *           "name": "Global"
         *         },
         *         {
         *           "name": "Shop Production"
         *         }
         *       ]
         *     }
         */
        TeamEnvironmentsUpdateAO: {
            /**
             * @description Environments that should be updated to this team
             * @example [
             *       {
             *         "name": "Global"
             *       },
             *       {
             *         "name": "Shop Production"
             *       }
             *     ]
             */
            environments: components["schemas"]["TeamEnvironmentAO"][];
        };
        /**
         * @description List of members that are part of this team.
         * @example {
         *       "members": [
         *         {
         *           "username": "13av2737-b318-4048-a79d-4789d645bc31",
         *           "role": "OWNER"
         *         },
         *         {
         *           "username": "google-oauth2|931412422966030837225",
         *           "role": "MEMBER"
         *         }
         *       ]
         *     }
         */
        TeamMembersAO: {
            /**
             * @description Members that are associated to this team
             * @example [
             *       {
             *         "username": "13av2737-b318-4048-a79d-4789d645bc31",
             *         "name": "jane doe",
             *         "email": "jane.doe@steadybit.com",
             *         "role": "OWNER"
             *       },
             *       {
             *         "username": "google-oauth2|931412422966030837225",
             *         "name": "john smith",
             *         "email": "jane.smith@steadybit.com",
             *         "role": "MEMBER"
             *       }
             *     ]
             */
            members: components["schemas"]["MemberAO"][];
        };
        /**
         * @description Team members that should be removed from a given team. You can specify the user to be removed via the internal identifier `usernames` or via the user's `emails`. If you specify both, both set of users will be removed.
         * @example {
         *       "emails": [
         *         "jane.doe@example.com",
         *         "javier.rodriguez@example.com"
         *       ],
         *       "usernames": [
         *         "13av2737-b318-4048-a79d-4789d645bc31"
         *       ]
         *     }
         */
        TeamMembersRemoveAO: {
            emails?: string[];
            usernames?: string[];
        };
        /**
         * @description Update request to change the members of a specific team. Specify either username, being a Steadybit user id, or the email address of the user.
         * @example {
         *       "members": [
         *         {
         *           "email": "jane.doe@example.com",
         *           "role": "OWNER"
         *         },
         *         {
         *           "email": "javier.rodriguez@example.com",
         *           "role": "MEMBER"
         *         },
         *         {
         *           "username": "auth0|1va2g15afc84590069cd53c3",
         *           "role": "MEMBER"
         *         }
         *       ]
         *     }
         */
        TeamMembersUpdateAO: {
            /**
             * @description Members that should be updated to this team
             * @example {
             *       "email": "jane.doe@example.com",
             *       "role": "OWNER"
             *     }
             */
            members: components["schemas"]["MemberUpdateAO"][];
        };
        /**
         * @description List of teams.
         * @example {
         *       "teams": [
         *         {
         *           "id": "714b0180-8abc-4d30-8acb-6aa024e3065f",
         *           "key": "ADM",
         *           "name": "Administrators",
         *           "version": 1,
         *           "logoId": "1",
         *           "logoColor": "cyanDark",
         *           "allowedActions": [
         *             "com.steadybit.extension_host.host.stress-cpu"
         *           ],
         *           "allowedEnvironments": [
         *             "Global"
         *           ],
         *           "members": [
         *             {
         *               "username": "auth0|11a9315afc84590069cd53b2",
         *               "role": "OWNER"
         *             }
         *           ]
         *         }
         *       ]
         *     }
         */
        TeamSummariesAO: {
            teams?: components["schemas"]["TeamAO"][];
        };
        /**
         * @description The tenant in which the event was performed. Only relevant in case you are using multiple tenants of the Steadybit platform.
         * @example {
         *       "key": "Demo",
         *       "name": "Demo Tenant"
         *     }
         */
        TenantAL: {
            key: string;
            name: string;
        };
        /** @description A named time series with date-value pairs. */
        TimeSeriesAO: {
            /** @description Name of the series, corresponding to the group label. */
            name?: string;
            /** @description Date-value pairs, each serialized as [date, count]. */
            values?: components["schemas"]["TimeSeriesValueAO"][];
        };
        /**
         * @description Time-series report data with metadata describing the query context.
         * @example {
         *       "from": "2026-01-01",
         *       "to": "2026-03-01",
         *       "rollup": "MONTHLY",
         *       "groupBy": "NONE",
         *       "series": [
         *         {
         *           "name": "users",
         *           "values": [
         *             [
         *               "2026-01-01",
         *               23
         *             ],
         *             [
         *               "2026-02-01",
         *               42
         *             ],
         *             [
         *               "2026-03-01",
         *               42
         *             ]
         *           ]
         *         }
         *       ]
         *     }
         */
        TimeSeriesReportAO: {
            /**
             * Format: date
             * @description Start date of the requested report range (inclusive).
             * @example 2026-01-01
             */
            from?: string;
            /**
             * @description The grouping dimension applied to the series.
             * @enum {string}
             */
            groupBy?: "NONE" | "STATE" | "TRIGGER" | "ACTION" | "CREATED_VIA" | "ORIGIN" | "ISSUES_FIXED" | "ISSUES_DISCOVERED" | "RISK_LEVEL" | "CATEGORY";
            /**
             * @description The time bucket granularity for report aggregation.
             * @enum {string}
             */
            rollup?: "MONTHLY" | "DAILY";
            /** @description The time-series data, one entry per group. */
            series?: components["schemas"]["TimeSeriesAO"][];
            /**
             * Format: date
             * @description End date of the requested report range (inclusive).
             * @example 2026-03-01
             */
            to?: string;
        };
        /** @description A date-value pair serialized as a two-element array [date, count]. */
        TimeSeriesValueAO: {
            /**
             * Format: date
             * @description Start date of the time bucket.
             * @example 2026-01-01
             */
            date?: string;
            /**
             * Format: int32
             * @description The count for this time bucket.
             */
            value?: number;
        };
        /**
         * @description Update the experiment with the given experiment design.
         * @example {
         *       "name": "Blackhole Hot-deals",
         *       "team": "ADM",
         *       "environment": "Global",
         *       "lanes": [
         *         {
         *           "steps": [
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "30s"
         *               },
         *               "actionType": "com.steadybit.extension_host.stress-cpu",
         *               "radius": {
         *                 "targetType": "com.steadybit.extension_container.container",
         *                 "predicate": {
         *                   "operator": "AND",
         *                   "predicates": [
         *                     {
         *                       "key": "k8s.deployment",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "hot-deals"
         *                       ]
         *                     }
         *                   ]
         *                 },
         *                 "query": null,
         *                 "percentage": 100
         *               }
         *             }
         *           ]
         *         }
         *       ],
         *       "properties": {
         *         "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
         *       }
         *     }
         */
        UpdateExperimentAO: {
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment?: string;
            /**
             * @description Variables that will be used when the experiment will be executed. Experiment variables will override existing environment variables. Each value is either a constant string, an array of constant strings, or a select expression object (`{"type":"select",...}`).
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal",
             *       "targetServices": [
             *         "gateway",
             *         "hot-deals",
             *         "fashion-bestseller"
             *       ],
             *       "httpEndpointZones": {
             *         "type": "select",
             *         "targetType": "com.steadybit.extension_container.container",
             *         "attribute": "aws.zone",
             *         "mode": "fixed",
             *         "count": 2
             *       }
             *     }
             */
            experimentVariables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * @description An optional external identifier used for create-or-update semantics.
             * @example 1234567
             */
            externalId?: string;
            /**
             * @deprecated
             * @description An optional external reference. Will be removed and is replaced by tags. If used with experiment creation, the value will be added as a tag.
             * @example INCIDENT-4711
             */
            externalReference?: string;
            /**
             * @description The hypothesis that is validated by the experiment
             * @example System is able to survive a latency in the network of 1500ms
             */
            hypothesis?: string;
            /**
             * @description The lanes (steps executed in parallel) in the experiment. Each lane consists of multiple steps that are executed sequential per lane.
             * @example [
             *       {
             *         "steps": [
             *           {
             *             "type": "action",
             *             "ignoreFailure": false,
             *             "parameters": {
             *               "duration": "30s"
             *             },
             *             "actionType": "com.steadybit.extension_host.stress-cpu",
             *             "radius": {
             *               "targetType": "com.steadybit.extension_container.container",
             *               "predicate": {
             *                 "operator": "AND",
             *                 "predicates": [
             *                   {
             *                     "key": "k8s.deployment",
             *                     "operator": "EQUALS",
             *                     "values": [
             *                       "hot-deals"
             *                     ]
             *                   }
             *                 ]
             *               },
             *               "query": null,
             *               "percentage": 100
             *             }
             *           }
             *         ]
             *       }
             *     ]
             */
            lanes: components["schemas"]["ExperimentLaneAO"][];
            /**
             * @description Name of the experiment to easily identify the experiment
             * @example Shop survives unavailability of hot-deals products
             */
            name: string;
            /**
             * @description The properties of the experiment
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * @description Team keys with which the experiment is shared with
             * @example [OPS, SHOP]
             */
            sharedTeams?: string[];
            /**
             * @description An optional set of tags you can use to search for.
             * @example [
             *       "myTag",
             *       "myOtherTag"
             *     ]
             */
            tags?: string[];
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            team: string;
        };
        /**
         * @description Experiment execution data that should be used only for that specific experiment execution and will not update the experiment design.
         * @example {
         *       "propertiesVersion": 1,
         *       "properties": {
         *         "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment execution!"
         *       },
         *       "propertiesOrder": [
         *         "EXAMPLE_CUSTOM_PROPERTY"
         *       ]
         *     }
         */
        UpdateExperimentExecutionPropertiesAO: {
            /**
             * @description The properties of the experiment execution
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment execution!"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * @description The order of the properties for this experiment execution. This may include global and experiment scoped assigned properties.
             * @example [
             *       "EXAMPLE_CUSTOM_PROPERTY"
             *     ]
             */
            propertiesOrder?: string[];
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            propertiesVersion?: number | null;
        };
        /**
         * @description Update an experiment based on an experiment template.
         * @example {
         *       "placeholders": [
         *         {
         *           "key": "CLUSTER",
         *           "value": "demo-cluster"
         *         },
         *         {
         *           "key": "BOOL",
         *           "value": true
         *         },
         *         {
         *           "key": "NUMBER",
         *           "value": 15
         *         },
         *         {
         *           "key": "KEYVALUE",
         *           "value": [
         *             {
         *               "key": "example-a",
         *               "value": "abc"
         *             },
         *             {
         *               "key": "example-b",
         *               "value": "123"
         *             }
         *           ]
         *         },
         *         {
         *           "key": "LIST",
         *           "value": [
         *             "entry1",
         *             "entry2",
         *             "entry3"
         *           ]
         *         },
         *         {
         *           "key": "FILE",
         *           "value": {
         *             "fileName": "example.txt",
         *             "data": "SGVsbG8gV29ybGQh"
         *           }
         *         }
         *       ],
         *     }
         */
        UpdateExperimentFromTemplateAO: {
            /** @description List of template placeholder values */
            placeholders?: components["schemas"]["ExperimentTemplatePlaceholderValueAO"][];
        };
        /**
         * @description Update or insert environment.
         * @example {
         *       "id": "2v1av42-e525-4c00-a13a-1ac32d170724",
         *       "name": "Global",
         *       "version": 0,
         *       "query": "aws.account=\"123\" OR aws.account=\"456\""
         *     }
         */
        UpsertEnvironmentAO: {
            /**
             * Format: uuid
             * @description Unique identifier of a environment
             */
            id?: string;
            /**
             * @description Name of the environment.
             * @example Global
             */
            name: string;
            predicate?: components["schemas"]["TargetPredicateAO"];
            /**
             * @description Alternative to `predicate`. If both `query` and `predicate` will be provided, `query` will override the `predicate`.
             * @example aws.account="123" OR aws.account="456"
             */
            query?: string | null;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number;
        };
        /**
         * @description An update or insert for an experiment schedule
         * @example {
         *       "experimentKey": "ADM-8",
         *       "cron": "30 * * * * ? *"
         *     }
         */
        UpsertExperimentScheduleAO: {
            /**
             * @description Should the experiment run if another experiment is running? Default is true.
             * @example true
             */
            allowParallel?: boolean;
            /**
             * @description Cron expression for the experiment schedule. Can't be used in combination with `startAt`.
             * @example 0 15 10 ? * *
             */
            cron?: string | null;
            /**
             * @description If `false`, the schedule is deactivated and no experiment will be executed. Default is true.
             * @example false
             */
            enabled?: boolean;
            /**
             * @description The experiment that should be scheduled.
             * @example ADM-123
             */
            experimentKey: string;
            /**
             * @description The unique identifier of the schedule. If not set, a new schedule will be created.
             * @example 01951394-727f-76a0-8675-c7519ebd0ff5
             */
            id?: string | null;
            /**
             * Format: date-time
             * @description Start date for a single execution. Can't be used in combination with `cron`.
             */
            startAt?: string | null;
            /**
             * @description Optional timezone for a experiment schedule. Can only be used with `cron`.
             * @example Europe/Berlin
             */
            timezone?: string | null;
            /**
             * @description Variables that will be used when the experiment will be executed. The variables will override existing environment or experiment variables. Each value is either a constant string, an array of constant strings, or a select expression object.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal"
             *     }
             */
            variables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
        };
        /**
         * @description Insert or update the experiment template in Steadybit. The `id` will be used to identify whether the template already exists and should be updated or newly inserted.
         * @example {
         *       "templateTitle": "HTTP Endpoint remains functional during Kubernetes Rollout Restart",
         *       "templateDescription": "Test if a given HTTP Endpoint remains funcitonal if a Kubernetes deployment is restarted.",
         *       "placeholders": [
         *         {
         *           "key": "HTTP_ENDPOINT",
         *           "name": "HTTP Endpoint",
         *           "description": "Which HTTP Endpoint should be checked during experiment execution?"
         *         },
         *         {
         *           "key": "DEPLOYMENT",
         *           "name": "Kubernetes Deployment",
         *           "description": "Which Kubernetes deployment do you want to restart?"
         *         },
         *         {
         *           "key": "CLUSTER",
         *           "name": "Kubernetes Cluster",
         *           "description": "In which Kubernetes cluster is the deployment deployed to?"
         *         },
         *         {
         *           "key": "NAMESPACE",
         *           "name": "Kubernetes Namespace",
         *           "description": "In which Kubernetes namespace is the deployment deployed to?"
         *         }
         *       ],
         *       "tags": [
         *         "Kubernetes"
         *       ],
         *       "lanes": [
         *         {
         *           "steps": [
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "60s",
         *                 "headers": [],
         *                 "method": "GET",
         *                 "successRate": "100",
         *                 "maxConcurrent": 5,
         *                 "followRedirects": false,
         *                 "readTimeout": "5s",
         *                 "connectTimeout": "5s",
         *                 "requestsPerSecond": 1,
         *                 "url": "[[HTTP_ENDPOINT]]",
         *                 "statusCode": "200-299"
         *               },
         *               "actionType": "com.steadybit.extension_http.check.periodically"
         *             }
         *           ]
         *         },
         *         {
         *           "steps": [
         *             {
         *               "type": "wait",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "10s"
         *               }
         *             },
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "wait": false
         *               },
         *               "actionType": "com.steadybit.extension_kubernetes.rollout-restart",
         *               "radius": {
         *                 "targetType": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *                 "predicate": {
         *                   "operator": "AND",
         *                   "predicates": [
         *                     {
         *                       "key": "k8s.cluster-name",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[CLUSTER]]"
         *                       ]
         *                     },
         *                     {
         *                       "key": "k8s.namespace",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[NAMESPACE]]"
         *                       ]
         *                     },
         *                     {
         *                       "key": "k8s.deployment",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[DEPLOYMENT]]"
         *                       ]
         *                     }
         *                   ]
         *                 },
         *                 "percentage": 50
         *               }
         *             },
         *             {
         *               "type": "action",
         *               "ignoreFailure": false,
         *               "parameters": {
         *                 "duration": "10m"
         *               },
         *               "actionType": "com.steadybit.extension_kubernetes.rollout-status",
         *               "radius": {
         *                 "targetType": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *                 "predicate": {
         *                   "operator": "AND",
         *                   "predicates": [
         *                     {
         *                       "key": "k8s.cluster-name",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[CLUSTER]]"
         *                       ]
         *                     },
         *                     {
         *                       "key": "k8s.namespace",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[NAMESPACE]]"
         *                       ]
         *                     },
         *                     {
         *                       "key": "k8s.deployment",
         *                       "operator": "EQUALS",
         *                       "values": [
         *                         "[[DEPLOYMENT]]"
         *                       ]
         *                     }
         *                   ]
         *                 },
         *                 "percentage": 50
         *               }
         *             }
         *           ]
         *         }
         *       ],
         *       "properties": {
         *         "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
         *       },
         *       "propertiesMetadata": [
         *         {
         *           "key": "EXAMPLE_CUSTOM_PROPERTY",
         *           "required": true,
         *           "editableInExecution": false
         *         }
         *       ]
         *     }
         */
        UpsertExperimentTemplateAO: {
            /**
             * @description Name of the experiment created by this template. If omitted, the name needs to be added when the template is used.
             * @example Shop survives unavailability of database
             */
            experimentName?: string | null;
            /**
             * @description Should the experiment template be hidden
             * @example false
             */
            hidden?: boolean;
            /**
             * @description The hypothesis that is validated by the experiment
             * @example System is able to survive a latency in the network of 1500ms
             */
            hypothesis?: string;
            /** Format: uuid */
            id?: string | null;
            /**
             * @description The lanes (steps executed in parallel) in the experiment template. Each lane consists of multiple steps that are executed sequential per lane.
             * @example [
             *       {
             *         "steps": [
             *           {
             *             "type": "action",
             *             "ignoreFailure": false,
             *             "parameters": {
             *               "duration": "30s"
             *             },
             *             "actionType": "com.steadybit.extension_host.stress-cpu",
             *             "radius": {
             *               "targetType": "com.steadybit.extension_container.container",
             *               "predicate": {
             *                 "operator": "AND",
             *                 "predicates": [
             *                   {
             *                     "key": "k8s.deployment",
             *                     "operator": "EQUALS",
             *                     "values": [
             *                       "hot-deals"
             *                     ]
             *                   }
             *                 ]
             *               },
             *               "query": null,
             *               "percentage": 100
             *             }
             *           }
             *         ]
             *       }
             *     ]
             */
            lanes: components["schemas"]["ExperimentLaneAO"][];
            /** @description A list of placeholders used in this experiment template. */
            placeholders?: components["schemas"]["ExperimentTemplatePlaceholderAO"][];
            /**
             * @description The properties of the experiment
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "I like this experiment!"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * @description Metadata for properties used in this template.
             * @example [
             *       {
             *         "key": "EXAMPLE_CUSTOM_PROPERTY",
             *         "required": true,
             *         "editableInExecution": false
             *       }
             *     ]
             */
            propertiesMetadata?: components["schemas"]["PropertyMetadataAO"][];
            /** @description A list of tags for this experiment template. (Up to 5) */
            tags?: string[];
            /** @description A brief description what the template is doing. */
            templateDescription: string;
            /**
             * @description The title of the template
             * @example Shop survives unavailability of database
             */
            templateTitle: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        UpsertHubAO: {
            /**
             * @description Website address of the the hub
             * @example https://hub.steadybit.com/
             */
            hubLink?: string | null;
            /** @description Name of the hub */
            hubName: string;
            /** Format: uuid */
            id?: string | null;
            /**
             * @description HTTP address of the the hub's repository
             * @example https://github.com/steadybit/reliability-hub-db
             */
            repositoryUrl: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /** @description Create or update a saved view of the explorer landscape. */
        UpsertLandscapeViewAO: {
            colorBy?: components["schemas"]["LandscapeViewColorByAO"];
            /**
             * @description Description of the saved view.
             * @example All shop workloads grouped by namespace.
             */
            description?: string;
            /**
             * @description Name of the environment the view is scoped to.
             * @example Global
             */
            environment?: string;
            /**
             * @description Explorer filter query narrowing the targets shown on the landscape.
             * @example k8s.namespace="shop"
             */
            filterQuery?: string;
            /** @description Ordered list of group-by dimensions the targets are grouped by, each with its own advanced configuration. */
            groupBy?: components["schemas"]["LandscapeViewGroupByAO"][];
            /**
             * @description Title of the saved view.
             * @example Kubernetes by namespace
             */
            name?: string;
            /**
             * @description Whether reliability advice is shown on the landscape.
             * @example false
             */
            showAdvice?: boolean;
            /**
             * @description Attribute key the size of a target is derived from.
             * @example k8s.container.cpu.limit
             */
            sizeBy?: string;
            /**
             * @description Key of the team the saved view belongs to.
             * @example ADM
             */
            team: string;
        };
        /**
         * @description A property association upsert.
         * @example {
         *       "key": "RESULT_COLOR",
         *       "editableInExecution": true,
         *       "required": true
         *     }
         */
        UpsertPropertyAssociationAO: {
            /**
             * @description Always defined to either `EXPERIMENT` for experiment design or run related associations or `SERVICE` for service-associations. Only for the former, an `experimentKey` can be defined and only for the latter, a `serviceId` can be defined
             * @default EXPERIMENT
             * @example EXPERIMENT
             * @enum {string}
             */
            associationType: "EXPERIMENT" | "SERVICE";
            /**
             * @description Is the property editable in the execution view. Only used when `associationType` is set to `EXPERIMENT`.
             * @example true
             */
            editableInExecution?: boolean;
            /**
             * @description The key of the associated experiment. When `associationType` is set to `EXPERIMENT` and `experimentKey` is `null`, it is associated to ALL experiment designs. Can't be changed during updates.
             * @example EXP-1
             */
            experimentKey?: string | null;
            /**
             * Format: uuid
             * @description Id of an existing Property-Association. A new association will be created if no id is provided or no matching association could be found
             */
            id?: string | null;
            /**
             * @description The key of the property definition
             * @example RESULT_COLOR
             */
            key: string;
            /**
             * @description Is the value required?
             * @example true
             */
            required?: boolean;
            /**
             * Format: uuid
             * @description The serviceId of the associated service. When `associationType` is set to `SERVICE` and `serviceId` is `null`, it is associated to ALL services. Can't be changed during updates.
             * @example 3308b47d-5c1f-4f08-a25b-a18fc10f8a56
             */
            serviceId?: string | null;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /**
         * @description A property association upsert.
         * @example {
         *       "key": "RESULT_COLOR",
         *       "label": "Result Color",
         *       "description": "How would you describe the result of your experiment, thinking in beautiful colors?",
         *       "dataType": "ENUM",
         *       "enumValues": [
         *         "RED",
         *         "GREEN",
         *         "BLUE"
         *       ]
         *     }
         */
        UpsertPropertyDefinitionAO: {
            /**
             * @description The data type of the property
             * @example STRING
             * @enum {string}
             */
            dataType: "STRING" | "STRING_LIST" | "ENUM" | "ENUM_LIST" | "NUMBER" | "NUMBER_LIST" | "MARKDOWN" | "BOOLEAN" | "DATE" | "LINK" | "LINK_LIST";
            /**
             * @description The text describing the property.
             * @example How would you describe the result of your experiment, thinking in beautiful colors?
             */
            description?: string | null;
            /**
             * @description Valid values if the dataType `ENUM` is used
             * @example [
             *       "RED",
             *       "GREEN",
             *       "BLUE"
             *     ]
             */
            enumValues?: string[];
            /**
             * @description The unique key of the property definition
             * @example RESULT_COLOR
             */
            key: string;
            /**
             * @description The label shown in the ui for this property
             * @example Result color
             */
            label: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        UpsertProvidedExperimentRequestAO: {
            /**
             * @description Update only - the experiment that should be updated
             * @example ADM-18
             */
            experimentKey?: string | null;
            /** @description List of template placeholder values */
            placeholders?: components["schemas"]["ExperimentTemplatePlaceholderValueAO"][] | null;
            /**
             * Format: uuid
             * @description The templateId that should be used for the provided experiment (needs to be included in the used service profile)
             * @example e1c22d74-a48b-4661-ab56-4e90c584c4e0
             */
            templateId: string;
        };
        /**
         * @example {
         *       "name": "shopping-service",
         *       "environment": "Global",
         *       "team": "ADM",
         *       "query": "aws.account=\"123\" OR aws.account=\"456\"",
         *       "validations": [
         *         {
         *           "type": "action",
         *           "parameters": {
         *             "url": "https://my-service/health",
         *             "method": "GET"
         *           },
         *           "actionType": "com.steadybit.extension_http.check.periodically"
         *         }
         *       ],
         *       "serviceProfile": "Steadybit Starter",
         *       "variables": {
         *         "httpEndpoint": "http://prod.shop.products.internal",
         *         "targets": {
         *           "type": "select",
         *           "targetType": "com.steadybit.extension_kubernetes.kubernetes-deployment",
         *           "attribute": "k8s.deployment",
         *           "filter": "k8s.namespace=\"shop\"",
         *           "mode": "fixed",
         *           "count": 1
         *         }
         *       }
         *     }
         */
        UpsertServiceAO: {
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment: string;
            /**
             * Format: uuid
             * @description The unique id of the service, will be created if not provided
             */
            id?: string | null;
            /**
             * @description Color scheme of the logo used to identify the service in the Platform UI
             * @default blue
             * @example orangeLight
             */
            logoColor: string;
            /**
             * @description Identifier of the logo used to identify the service in the Platform UI
             * @default service
             * @example service-router
             */
            logoId: string;
            /**
             * @description The name of the service
             * @example calculator-service
             */
            name: string;
            /**
             * @description The properties of the service
             * @example {
             *       "EXAMPLE_CUSTOM_PROPERTY": "I like this service!"
             *     }
             */
            properties?: {
                [key: string]: unknown;
            };
            /**
             * @description Query-Language predicate, specifies the targets belonging to this Service
             * @example aws.account="123" OR aws.account="456"
             */
            query: string;
            /**
             * @description Name of the service profile that should be used for this service
             * @example Steadybit provided
             */
            serviceProfile: string;
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            team: string;
            /** @description List of validations to be executed against the service */
            validations: components["schemas"]["ExperimentStepActionAO"][];
            /**
             * @description Variables owned by the service. Each value is either a constant string, an array of constant strings, or a select expression object. A select-expression value **must set `scope`** (`service` or `environment`) — the request is rejected otherwise. On `POST /api/services` (upsert): omitting this field leaves existing variables untouched, an empty object removes all of them.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal",
             *       "targetServices": [
             *         "gateway",
             *         "hot-deals",
             *         "fashion-bestseller"
             *       ]
             *     }
             */
            variables?: {
                [key: string]: components["schemas"]["VariableExpressionAO"];
            };
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /**
         * @description Request to create or update a service profile
         * @example {
         *       "name": "My Custom Templates",
         *       "origin": "CUSTOM",
         *       "templates": {
         *         "availability": [
         *           "550e8400-e29b-41d4-a716-446655440000"
         *         ],
         *         "latency": [
         *           "550e8400-e29b-41d4-a716-446655440002"
         *         ]
         *       }
         *     }
         */
        UpsertServiceProfileAO: {
            /**
             * Format: uuid
             * @description The unique id of the profile. Will be created if not provided.
             */
            id?: string | null;
            /**
             * @description The name of the profile
             * @example Default Resilience Tests
             */
            name: string;
            /**
             * @description Origin of a service profile
             * @example CUSTOM
             * @enum {string}
             */
            origin: "PROVIDED" | "CUSTOM";
            /** @description Template entries in this profile */
            templates: components["schemas"]["ServiceProfileCategoryAO"][];
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /**
         * @description Insert or update the team in Steadybit. The `key` will be used to identify whether the team exists already and should be updated or newly inserted.
         * @example {
         *       "id": "71ab0180-8abc-4d30-8acb-6aa024e3065f",
         *       "key": "ADM",
         *       "name": "Administrators",
         *       "version": 1,
         *       "logoId": "1",
         *       "logoColor": "cyanDark",
         *       "allowedActions": [
         *         "com.steadybit.extension_host.host.stress-cpu"
         *       ],
         *       "allowedEnvironments": [
         *         "Global"
         *       ],
         *       "members": [
         *         {
         *           "email": "jane.doe@example.com",
         *           "role": "OWNER"
         *         },
         *         {
         *           "email": "javier.rodriguez@example.com",
         *           "role": "MEMBER"
         *         },
         *         {
         *           "username": "auth0|1va2g15afc84590069cd53c3",
         *           "role": "MEMBER"
         *         }
         *       ]
         *     }
         */
        UpsertTeamAO: {
            /**
             * @description Set of allowed actions that can be used in an experiment of this team
             * @example [
             *       "com.steadybit.extension_host.host.stress-cpu"
             *     ]
             */
            allowedActions: string[];
            /**
             * @description Set of allowed environments, identified via name
             * @example [
             *       "Global"
             *     ]
             */
            allowedEnvironments: string[];
            /** @description An optional description of a team */
            description?: string;
            /** Format: uuid */
            id?: string;
            /**
             * @description Unique identifier of a team
             * @example ADM
             */
            key: string;
            /**
             * @description Color scheme of the logo used to identify the team in the Platform UI
             * @example cyanDark
             */
            logoColor?: string;
            /**
             * @description Identifier of the logo used to identify the team in the Platform UI
             * @example 1
             */
            logoId?: string;
            /**
             * @description How a team or team membership is managed
             * @example MANUAL
             * @enum {string}
             */
            managedBy?: "MANUAL" | "OIDC" | "LDAP";
            /**
             * @description Members that should be added to this team
             * @example {
             *       "email": "jane.doe@example.com",
             *       "role": "OWNER"
             *     }
             */
            members?: components["schemas"]["MemberUpdateAO"][];
            /**
             * @description Name of a team
             * @example ADMIN
             */
            name: string;
            /**
             * Format: int32
             * @description Version for optimistic locking (optional in the API)
             * @example 1
             */
            version?: number | null;
        };
        /**
         * @description A user has performed the logged event e.g. via UI
         * @example {
         *       "name": "Jane Doe",
         *       "role": "ADMIN",
         *       "email": "example@example.com",
         *       "username": "1ava2afg-xju33-4c6a-9451-2854584c15be",
         *       "principalType": "USER"
         *     }
         */
        UserPrincipalAL: {
            /**
             * @description E-mail of the user, unique within Steadybit
             * @example example@example.com
             */
            email: string;
            /**
             * @description Name of the user
             * @example Jane Doe
             */
            name: string;
            /**
             * @description Principal type for user based principal
             * @example USER
             * @enum {string}
             */
            principalType: "USER" | "ACCESS_TOKEN" | "BATCH_JOB";
            /**
             * @description Role of the user in the platform
             * @example ADMIN
             * @enum {string|null}
             */
            role?: "ADMIN" | "USER" | "SUPPORT" | null;
            /**
             * @description Username of the user, internal identifier of Steadybit
             * @example 13av2737-b318-4048-a79d-4789d645bc31
             */
            username: string;
        };
        /**
         * @description The user that canceled the experiment execution, only present if the execution was canceled
         * @example {
         *       "username": "ag1hb7ap-d299-47ab-998f-c2a53b433820",
         *       "name": "Max Mustermann",
         *       "email": "max@steadybit.com"
         *     }
         */
        UserSummaryAO: {
            email?: string | null;
            /**
             * @description Name of the user
             * @example Jane Doe
             */
            name?: string;
            pictureUrl?: string | null;
            /**
             * @description Username of the user, internal identifier of Steadybit
             * @example 13av2737-b318-4048-a79d-4789d645bc31
             */
            username: string;
        };
        /** @description A variable value: a constant string (≤5000 chars), an array of constant strings, or a select expression object. */
        VariableExpressionAO: string | string[] | components["schemas"]["SelectExpressionAO"];
        VariableValueAO: {
            type: "VariableValueAO";
        } & (Omit<components["schemas"]["ComparableValueAO"], "type"> & {
            value: string;
        });
        /**
         * @description Webhook payload containing more event-specific information.
         * @example {
         *       "event": "killswitch.disengaged",
         *       "time": "2023-01-01T09:15:00.000000Z",
         *       "killswitch": {
         *         "engagedBy": "admin",
         *         "engaged": "2023-01-01T09:00:00.000000Z",
         *         "disengagedBy": "admin",
         *         "disengaged": "2023-01-01T09:15:00.000000Z"
         *       }
         *     }
         */
        WebhookPayloadAO: {
            /**
             * @description The event that caused the webhook to be triggered
             * @example killswitch.disengaged
             */
            event?: string;
            execution?: components["schemas"]["WebhookPayloadExecutionAO"];
            /**
             * Format: uuid
             * @description In case the webhook was triggered by an experiment execution step event, this contains the actual step's id having triggered the webhook
             * @example 32009f6e-ff90-47a5-8daa-80f4bb7ae591
             */
            executionStepId?: string;
            killswitch?: components["schemas"]["WebhookPayloadKillswitchAO"];
            /**
             * Format: date-time
             * @description The timestamp at which the event was fired
             * @example 2023-01-01T09:15:00Z
             */
            time?: string;
        };
        /**
         * @description The execution of a single experiment.
         * @example {
         *       "id": 1,
         *       "experimentKey": "ADM-1",
         *       "teamKey": "ADM",
         *       "environment": "32009f6e-ff90-47a5-8daa-80f4bb7ae591",
         *       "name": "My first experiment",
         *       "created": "2023-01-01T08:00:00.000000Z",
         *       "createdVia": "UI",
         *       "experimentVersion": "5",
         *       "state": "CREATED",
         *       "steps": [
         *         {
         *           "ignoreFailure": false,
         *           "parameters": {
         *             "duration": "10s"
         *           }
         *         },
         *         {
         *           "predecessorId": "40b0f797-912d-4256-8887-1553561962a9",
         *           "ignoreFailure": false,
         *           "parameters": {
         *             "cpuLoad": 100,
         *             "workers": 0,
         *             "duration": "30s"
         *           },
         *           "actionId": "com.steadybit.extension_container.stress_cpu",
         *           "actionKind": "ATTACK",
         *           "radius": {
         *             "targetType": "com.steadybit.extension_container.container",
         *             "percentage": 50,
         *             "predicate": {
         *               "operator": "AND",
         *               "predicates": [
         *                 {
         *                   "key": "container.host/name",
         *                   "operator": "EQUALS",
         *                   "values": [
         *                     "docker-desktop/minikube"
         *                   ]
         *                 }
         *               ]
         *             }
         *           },
         *           "targetExecutions": [
         *             {
         *               "type": "com.steadybit.extension_container.container",
         *               "name": "docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea",
         *               "state": "COMPLETED",
         *               "attributes": [
         *                 {
         *                   "key": "container.port",
         *                   "value": "51152:2376"
         *                 },
         *                 {
         *                   "key": "container.engine",
         *                   "value": "docker"
         *                 },
         *                 {
         *                   "key": "container.host/name",
         *                   "value": "docker-desktop/minikube"
         *                 },
         *                 {
         *                   "key": "container.host",
         *                   "value": "docker-desktop"
         *                 }
         *               ]
         *             }
         *           ],
         *           "totalTargetCount": 1
         *         }
         *       ]
         *     }
         */
        WebhookPayloadExecutionAO: {
            canceledBy?: components["schemas"]["UserSummaryAO"];
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution was created
             * @example 2023-01-01T09:00:01Z
             */
            created?: string;
            createdBy?: components["schemas"]["UserSummaryAO"];
            /**
             * @description The creation trigger that caused this experiment execution to be started
             * @example UI
             * @enum {string}
             */
            createdVia?: "API" | "CLI" | "UI" | "SCHEDULE" | "MCP" | "SUITE";
            /**
             * Format: date-time
             * @description Timestamp when the experiment ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * @description The name of the environment to be used
             * @example Global
             */
            environment?: string;
            /**
             * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
             * @example ADM-2
             */
            experimentKey?: string;
            /**
             * Format: int32
             * @description Experiment design version which can be used to identify changes between experiment runs
             * @example 5
             */
            experimentVersion?: number;
            /**
             * @description The hypothesis that is validated by the experiment
             * @example System is able to survive a latency in the network of 1500ms
             */
            hypothesis?: string;
            /**
             * Format: int32
             * @description Unique experiment execution id that identifies this specific experiment execution
             * @example 1523
             */
            id?: number;
            /**
             * @description Name of the experiment to easily identify the experiment
             * @example Shop survives unavailability of hot-deals products
             */
            name?: string;
            overrides?: components["schemas"]["WebhookPayloadExecutionOverridesAO"];
            /**
             * @description Reason in case the experiment execution failed or errored
             * @example Action error
             */
            reason?: string;
            /**
             * @description Additional detail for the reason in case the experiment execution failed or errored
             * @example Couldn't read state of container...
             */
            reasonDetails?: string;
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution was requested
             * @example 2023-01-01T09:00:00Z
             */
            requested?: string;
            /**
             * Format: date-time
             * @description Timestamp when the experiment execution started running, after preparation/preflight
             * @example 2023-01-01T09:00:05Z
             */
            started?: string;
            /**
             * @description Current state of the experiment (e.g. CREATED, RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
            /**
             * @description The steps that are executed in parallel or sequence in the experiment.
             * @example [
             *       {
             *         "ignoreFailure": false,
             *         "parameters": {
             *           "duration": "10s"
             *         }
             *       },
             *       {
             *         "predecessorId": "40b0f797-912d-4256-8887-1553561962a9",
             *         "ignoreFailure": false,
             *         "parameters": {
             *           "cpuLoad": 100,
             *           "workers": 0,
             *           "duration": "30s"
             *         },
             *         "actionId": "com.steadybit.extension_container.stress_cpu",
             *         "actionKind": "ATTACK",
             *         "radius": {
             *           "targetType": "com.steadybit.extension_container.container",
             *           "percentage": 50,
             *           "predicate": {
             *             "operator": "AND",
             *             "predicates": [
             *               {
             *                 "key": "container.host/name",
             *                 "operator": "EQUALS",
             *                 "values": [
             *                   "docker-desktop/minikube"
             *                 ]
             *               }
             *             ]
             *           }
             *         },
             *         "targetExecutions": [
             *           {
             *             "type": "com.steadybit.extension_container.container",
             *             "name": "docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea",
             *             "state": "COMPLETED",
             *             "attributes": [
             *               {
             *                 "key": "container.port",
             *                 "value": "51152:2376"
             *               },
             *               {
             *                 "key": "container.engine",
             *                 "value": "docker"
             *               },
             *               {
             *                 "key": "container.host/name",
             *                 "value": "docker-desktop/minikube"
             *               },
             *               {
             *                 "key": "container.host",
             *                 "value": "docker-desktop"
             *               }
             *             ]
             *           }
             *         ],
             *         "totalTargetCount": 1
             *       }
             *     ]
             */
            steps?: components["schemas"]["AbstractWebhookPayloadExecutionStepAO"][];
            /**
             * @description Tags of the experiment at the time the execution was requested
             * @example [
             *       "resilience",
             *       "shop"
             *     ]
             */
            tags?: string[];
            /**
             * @description The key of the team to be used
             * @example ADM
             */
            teamKey?: string;
            /**
             * @description The variables resolved for this specific execution, keyed by name. Each entry carries the resolved value(s) and the tier the winning value originated from (ENVIRONMENT, SERVICE, EXPERIMENT, SCHEDULE, EXECUTION). A single-value variable's value is a string, a multi-value variable's value is an array of strings. Empty until the execution starts, as dynamic values are resolved once at run start and then stay stable for the whole run.
             * @example {
             *       "httpEndpoint": {
             *         "value": "http://shop.products.internal",
             *         "origin": "EXECUTION"
             *       }
             *     }
             */
            variables?: {
                [key: string]: components["schemas"]["ExperimentExecutionVariableAO"];
            };
        };
        /**
         * @deprecated
         * @description Overrides that are used for a single experiment execution and not saved into the experiment design
         * @example {
         *       "environment": "Shop-DEV",
         *       "variables": {
         *         "httpEndpoint": "http://dev.shop.products.internal"
         *       }
         *     }
         */
        WebhookPayloadExecutionOverridesAO: {
            /**
             * @description The environment in which the experiment should be executed instead of the environment specified in the experiment design
             * @example Shop-DEV
             */
            environment?: string;
            /**
             * @description A set of environment variables that should only be used for that specific experiment execution. Constant values are passed literally; dynamic (select expression) values are rendered as a human-readable summary, e.g. `select 3 of host.name on com.steadybit.extension_host.host`.
             * @example {
             *       "httpEndpoint": "http://dev.shop.products.internal"
             *     }
             */
            variables?: {
                [key: string]: string;
            };
        };
        /**
         * @description An action-step that is executed as part of an experiment.
         * @example {
         *       "predecessorId": "40b0f797-912d-4256-8887-1553561962a9",
         *       "ignoreFailure": false,
         *       "parameters": {
         *         "cpuLoad": 100,
         *         "workers": 0,
         *         "duration": "30s"
         *       },
         *       "actionId": "com.steadybit.extension_container.stress_cpu",
         *       "actionKind": "ATTACK",
         *       "radius": {
         *         "targetType": "com.steadybit.extension_container.container",
         *         "percentage": 50,
         *         "predicate": {
         *           "operator": "AND",
         *           "predicates": [
         *             {
         *               "key": "container.host/name",
         *               "operator": "EQUALS",
         *               "values": [
         *                 "docker-desktop/minikube"
         *               ]
         *             }
         *           ]
         *         }
         *       },
         *       "targetExecutions": [
         *         {
         *           "type": "com.steadybit.extension_container.container",
         *           "name": "docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea",
         *           "state": "COMPLETED",
         *           "attributes": [
         *             {
         *               "key": "container.port",
         *               "value": "51152:2376"
         *             },
         *             {
         *               "key": "container.engine",
         *               "value": "docker"
         *             },
         *             {
         *               "key": "container.host/name",
         *               "value": "docker-desktop/minikube"
         *             },
         *             {
         *               "key": "container.host",
         *               "value": "docker-desktop"
         *             }
         *           ]
         *         }
         *       ],
         *       "totalTargetCount": 1
         *     }
         */
        WebhookPayloadExecutionStepActionAO: {
            /**
             * @description Unique identifier of the action that is executed in this step
             * @example com.steadybit.extension_container.stress_cpu
             */
            actionId?: string;
            /**
             * @description Kind of the action (e.g. attack, check, loadtest)
             * @example ATTACK
             * @enum {string}
             */
            actionKind?: "ATTACK" | "CHECK" | "LOAD_TEST" | "OTHER" | "BASIC";
            /**
             * @description Custom label assigned during experiment design to express the intention of this step
             * @example Container 'xyz' can not be reached
             */
            customLabel?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * Format: uuid
             * @description Unique identifier of this step execution
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            id?: string;
            /**
             * @description Whether the experiment should fail/error immediately in case this step fails/errors.
             * @example false
             */
            ignoreFailure?: boolean;
            /**
             * @description Step-specific parameters of the experiment step configuration
             * @example {
             *       "duration": "10s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            /**
             * Format: uuid
             * @description Unique identifier of the step execution that precedes this step, null if it is the first step of a lane
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            predecessorId?: string;
            radius?: components["schemas"]["WebhookPayloadExecutionStepActionBlastRadiusAO"];
            /**
             * @description Reason in case this experiment step execution failed or errored
             * @example Couldn't read state of container...
             */
            reason?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step was started
             * @example 2023-01-01T09:00:00Z
             */
            started?: string;
            /**
             * @description Current state of this step in the experiment (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
            /**
             * @description List of targets that are expected to be effected by this action. This list may change in case targets aren't available at the specific time of execution
             * @example [
             *       {
             *         "type": "com.steadybit.extension_container.container",
             *         "name": "docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea",
             *         "state": "COMPLETED",
             *         "attributes": [
             *           {
             *             "key": "container.port",
             *             "value": "51152:2376"
             *           },
             *           {
             *             "key": "container.engine",
             *             "value": "docker"
             *           },
             *           {
             *             "key": "container.host/name",
             *             "value": "docker-desktop/minikube"
             *           },
             *           {
             *             "key": "container.host",
             *             "value": "docker-desktop"
             *           }
             *         ]
             *       }
             *     ]
             */
            targetExecutions?: components["schemas"]["WebhookPayloadExecutionStepActionTargetAO"][];
            /**
             * Format: int64
             * @description Amount of targets that are effect int total
             * @example 23
             */
            totalTargetCount?: number;
        };
        /**
         * @description Blast radius that is applied to define the set of targets as well as an optional random subset
         * @example {
         *       "targetType": "com.steadybit.extension_container.container",
         *       "percentage": 50,
         *       "predicate": {
         *         "operator": "AND",
         *         "predicates": [
         *           {
         *             "key": "container.host/name",
         *             "operator": "EQUALS",
         *             "values": [
         *               "docker-desktop/minikube"
         *             ]
         *           }
         *         ]
         *       }
         *     }
         */
        WebhookPayloadExecutionStepActionBlastRadiusAO: {
            /**
             * Format: int32
             * @description In case a fixed number of as subset of specified targets should be effected
             * @example 2
             */
            maximum?: number;
            /**
             * Format: int32
             * @description In case a percentage subset of the specified targets should be effected
             * @example 40
             */
            percentage?: number;
            predicate?: components["schemas"]["TargetPredicateAO"];
            /**
             * @description Target type that is effected by that action
             * @example container
             */
            targetType?: string;
        };
        /**
         * @description A target that is expected to be effected by this action.
         * @example {
         *       "type": "com.steadybit.extension_container.container",
         *       "name": "docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea",
         *       "state": "COMPLETED",
         *       "attributes": [
         *         {
         *           "key": "container.port",
         *           "value": "51152:2376"
         *         },
         *         {
         *           "key": "container.engine",
         *           "value": "docker"
         *         },
         *         {
         *           "key": "container.host/name",
         *           "value": "docker-desktop/minikube"
         *         },
         *         {
         *           "key": "container.host",
         *           "value": "docker-desktop"
         *         }
         *       ]
         *     }
         */
        WebhookPayloadExecutionStepActionTargetAO: {
            /**
             * @description A set of attributes that have been discovered for this target. A key may be associated multiple time to a single target.
             * @example [
             *       {
             *         "key": "container.port",
             *         "value": "51152:2376"
             *       },
             *       {
             *         "key": "container.engine",
             *         "value": "docker"
             *       }
             *     ]
             */
            attributes?: components["schemas"]["Attribute"][];
            /**
             * @description Identifier of the target that is expected to be effected
             * @example docker://1f769d01b9c5cd29bb302ca40157274f38798104208117b3310825ba676883ea
             */
            name?: string;
            /**
             * @description Type of the target that is expected to be effected
             * @example container
             */
            type?: string;
        };
        /**
         * @description A service validation step that is executed as part of an experiment.
         * @example {
         *       "stepType": "SERVICE-VALIDATION",
         *       "id": "40b0f797-912d-4256-8887-1553561962a9",
         *       "state": "COMPLETED",
         *       "started": "2025-06-18T08:32:01.850479Z",
         *       "ended": "2025-06-18T08:32:11.886043Z",
         *       "predecessorId": null,
         *       "ignoreFailure": false,
         *       "parameters": {
         *         "duration": "60s"
         *       },
         *       "serviceId": "cc06f132-0694-4ffa-aee2-13d8fafa3a8b",
         *       "validations": []
         *     }
         */
        WebhookPayloadExecutionStepServiceValidationAO: {
            /**
             * @description Custom label assigned during experiment design to express the intention of this step
             * @example Container 'xyz' can not be reached
             */
            customLabel?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * Format: uuid
             * @description Unique identifier of this step execution
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            id?: string;
            /**
             * @description Whether the experiment should fail/error immediately in case this step fails/errors.
             * @example false
             */
            ignoreFailure?: boolean;
            /**
             * @description Step-specific parameters of the experiment step configuration
             * @example {
             *       "duration": "10s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            /**
             * Format: uuid
             * @description Unique identifier of the step execution that precedes this step, null if it is the first step of a lane
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            predecessorId?: string;
            /**
             * @description Reason in case this experiment step execution failed or errored
             * @example Couldn't read state of container...
             */
            reason?: string;
            /**
             * Format: uuid
             * @description Unique identifier of the service.
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            serviceId?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step was started
             * @example 2023-01-01T09:00:00Z
             */
            started?: string;
            /**
             * @description Current state of this step in the experiment (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
            /** @description List of actions performed as part of this service validation step. */
            validations?: components["schemas"]["WebhookPayloadExecutionStepActionAO"][];
        };
        /**
         * @description A wait step that is executed as part of an experiment.
         * @example {
         *       "id": "40b0f797-912d-4256-8887-1553561962a9",
         *       "predecessorId": null,
         *       "ignoreFailure": false,
         *       "parameters": {
         *         "duration": "10s"
         *       }
         *     }
         */
        WebhookPayloadExecutionStepWaitAO: {
            /**
             * @description Custom label assigned during experiment design to express the intention of this step
             * @example Container 'xyz' can not be reached
             */
            customLabel?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step ended
             * @example 2023-01-01T09:00:00Z
             */
            ended?: string;
            /**
             * Format: uuid
             * @description Unique identifier of this step execution
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            id?: string;
            /**
             * @description Whether the experiment should fail/error immediately in case this step fails/errors.
             * @example false
             */
            ignoreFailure?: boolean;
            /**
             * @description Step-specific parameters of the experiment step configuration
             * @example {
             *       "duration": "10s"
             *     }
             */
            parameters?: {
                [key: string]: unknown;
            };
            /**
             * Format: uuid
             * @description Unique identifier of the step execution that precedes this step, null if it is the first step of a lane
             * @example 40b0f797-912d-4256-8887-1553561962a9
             */
            predecessorId?: string;
            /**
             * @description Reason in case this experiment step execution failed or errored
             * @example Couldn't read state of container...
             */
            reason?: string;
            /**
             * Format: date-time
             * @description Timestamp when this experiment step was started
             * @example 2023-01-01T09:00:00Z
             */
            started?: string;
            /**
             * @description Current state of this step in the experiment (e.g. RUNNING, FAILED, ERRORED, COMPLETED)
             * @example RUNNING
             */
            state?: string;
        };
        /**
         * @description Webhook payload performed for a killswitch related event.
         * @example {
         *       "engagedBy": "13av2737-b318-4048-a79d-4789d645bc31",
         *       "engaged": "2023-01-01T09:00:00.000000Z",
         *       "disengagedBy": "13av2737-b318-4048-a79d-4789d645bc31",
         *       "disengaged": "2023-01-01T09:15:00.000000Z"
         *     }
         */
        WebhookPayloadKillswitchAO: {
            /**
             * Format: date-time
             * @description Timestamp when the kill switch was disengaged
             * @example 2023-01-01T09:15:00Z
             */
            disengaged?: string;
            /**
             * @description Username of the user who has disengaged the kill switch
             * @example 13av2737-b318-4048-a79d-4789d645bc31
             */
            disengagedBy?: string;
            disengagedByDetails?: components["schemas"]["UserSummaryAO"];
            /**
             * Format: date-time
             * @description Timestamp when the kill switch was engaged
             * @example 2023-01-01T09:00:00Z
             */
            engaged?: string;
            /**
             * @description Username of the user who has engaged the kill switch
             * @example 13av2737-b318-4048-a79d-4789d645bc31
             */
            engagedBy?: string;
            engagedByDetails?: components["schemas"]["UserSummaryAO"];
        };
    };
    responses: never;
    parameters: never;
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export interface operations {
    getAccessTokens: {
        parameters: {
            query: {
                pageRequest: components["schemas"]["PageRequestAO"];
                team?: string;
                type?: "ADMIN" | "TEAM";
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOAccessTokensPageItemAO"];
                    "application/yaml": components["schemas"]["PagedResponseAOAccessTokensPageItemAO"];
                };
            };
        };
    };
    createAccessToken: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CreateAccessTokenRequestAO"];
            };
        };
        responses: {
            /** @description Token created, Response body contains the newly generated token.<br/> Make sure to save the generated token as you can't read it again afterwards for security-reasons. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CreateAccessTokenResponseAO"];
                    "application/yaml": components["schemas"]["CreateAccessTokenResponseAO"];
                };
            };
            /** @description Validation error, the access token was not generated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteAccessToken: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Token deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Token not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, the access token was not deleted */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getAccessTokens_1: {
        parameters: {
            query: {
                createdBy?: string;
                expired?: boolean;
                name?: string;
                pageRequest: components["schemas"]["PageRequestAO"];
                teams?: string[];
                type?: "ADMIN" | "TEAM" | "WILDCARD";
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOAccessTokensPageItemV2AO"];
                    "application/yaml": components["schemas"]["PagedResponseAOAccessTokensPageItemV2AO"];
                };
            };
        };
    };
    createAccessToken_1: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CreateAccessTokenRequestV2AO"];
            };
        };
        responses: {
            /** @description Token created. Make sure to save the generated token as you can't read it again afterwards for security-reasons. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CreateAccessTokenResponseV2AO"];
                    "application/yaml": components["schemas"]["CreateAccessTokenResponseV2AO"];
                };
            };
            /** @description Validation error, the access token was not generated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteAccessToken_1: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Token deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Token not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    recreateAccessToken: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["RecreateAccessTokenRequestV2AO"];
            };
        };
        responses: {
            /** @description Token recreated. Make sure to save the generated token as you can't read it again afterwards for security-reasons. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CreateAccessTokenResponseV2AO"];
                    "application/yaml": components["schemas"]["CreateAccessTokenResponseV2AO"];
                };
            };
            /** @description Token not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    findAllActions: {
        parameters: {
            query?: {
                /** @description The page number to retrieve. Starts from 0. */
                page?: number;
                /** @description The number of items to return per page. */
                size?: number;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Actions including their parameters, limited to the given page size. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ActionSummariesAO"];
                    "application/yaml": components["schemas"]["ActionSummariesAO"];
                };
            };
        };
    };
    getAction: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Action ID
                 * @example com.steadybit.extension_container.network_block_dns
                 */
                actionId: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Action found, response body contains the action details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ActionAO"];
                    "application/yaml": components["schemas"]["ActionAO"];
                };
            };
            /** @description Action with given `actionId` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getTargetAdviceSummary: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["GetAdviceApiRequestAO"];
            };
        };
        responses: {
            /** @description Matching targets with advice information, limited to 20 items. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["AdviceSummaryAO"];
                    "application/yaml": components["schemas"]["AdviceSummaryAO"];
                };
            };
        };
    };
    find: {
        parameters: {
            query?: {
                /**
                 * @description Starting point with the earliest time to be included.<br/>If neither `to` nor `from` is specified, it defaults to a 7 days date range from today.
                 * @example 2023-01-01T09:00:00Z
                 */
                from?: string | null;
                /**
                 * @description End point with the latest time to be included.<br/>If neither `to` nor `from` is specified, it defaults to a 7 days date range from today.
                 * @example 2023-01-31T09:00:00Z
                 */
                to?: string | null;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description List of all audit logs */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["AuditLogEntry"][];
                };
            };
        };
    };
    forwardToPlatform: {
        parameters: {
            query: {
                /**
                 * @deprecated
                 * @description External reference that identifies the experiment. This is used to identify whether an experiment was already created for that reference or not. Can be e.g. an incident or ticket identifier of pager duty or JIRA
                 * @example INCIDENT-312
                 */
                externalReference?: string;
                /**
                 * @description Tag that identifies the experiment. This is used to identify whether an experiment was already created for these tag or not. Can be e.g. an incident or ticket identifier of pager duty or JIRA
                 * @example INCIDENT,INCIDENT-312
                 */
                tag?: string;
                /**
                 * @description Key of the Steadybit tenant. You can get the key from the Platform URL or by asking the Steadybit team
                 * @example demo
                 */
                tenantKey: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Forwarding to create experiment or (a list of) experiments */
            307: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": Record<string, never>;
                    "application/yaml": Record<string, never>;
                };
            };
        };
    };
    getLinkedBadge: {
        parameters: {
            query: {
                /**
                 * @description Caption that is shown at the badge when no experiment exists in order to create a new experiment
                 * @example Create experiment
                 */
                createCaption?: string;
                /**
                 * @deprecated
                 * @description External reference that identifies the experiment. This is used to identify whether an experiment was already created for that reference or not. Can be e.g. an incident or ticket identifier of pager duty or JIRA
                 * @example INCIDENT-312
                 */
                externalReference?: string;
                /**
                 * @description Optional parameter in case you need to scale the svg image. Defaults to 1
                 * @example 2
                 */
                scale?: number;
                /**
                 * @description A tag that identifies the experiment. This is used to identify whether an experiment was already created having this tag or not. Can be e.g. an incident or ticket identifier of pager duty or JIRA
                 * @example INCIDENT,INCIDENT-312
                 */
                tag?: string;
                /**
                 * @description Key of the Steadybit tenant (only for SaaS customers). You can get the key from the Platform URL or by asking the Steadybit team
                 * @example demo
                 */
                tenantKey: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Returns a badge as a svg-image */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "image/svg+xml": string;
                };
            };
        };
    };
    getEnvironments: {
        parameters: {
            query?: {
                /**
                 * @description If set, only environments matching the search are returned. Matches the environment name or the name or key of a team the environment is assigned to.
                 * @example shop
                 */
                search?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all existing environments. Fetch a single environment by `id` to get more details for a team */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["EnvironmentSummariesAO"];
                    "application/yaml": components["schemas"]["EnvironmentSummariesAO"];
                };
            };
        };
    };
    upsertEnvironment: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertEnvironmentAO"];
            };
        };
        responses: {
            /** @description Environment updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["EnvironmentAO"];
                    "application/yaml": components["schemas"]["EnvironmentAO"];
                };
            };
            /** @description Environment created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["EnvironmentAO"];
                    "application/yaml": components["schemas"]["EnvironmentAO"];
                };
            };
            /** @description Validation error, environment was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getEnvironment: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a environment
                 * @example 2v1av42-e525-4c00-a13a-1ac32d170724
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Environment found, response body contains the environment details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["EnvironmentAO"];
                    "application/yaml": components["schemas"]["EnvironmentAO"];
                };
            };
            /** @description The given `id` is not a uuid */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Environment with given `id` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteEnvironment: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a environment
                 * @example 2v1av42-e525-4c00-a13a-1ac32d170724
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Environment deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Environment with given `id` was not found and thus not deleted */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getEnvironmentVariables: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a environment
                 * @example 2v1av42-e525-4c00-a13a-1ac32d170724
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Key-value-Map of environment variables associated to the environment. Each value is either a constant string, an array of constant strings, or a select expression object. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": string;
                    "application/yaml": string;
                };
            };
            /** @description Environment not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    updateEnvironmentVariables: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a environment
                 * @example 2v1av42-e525-4c00-a13a-1ac32d170724
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": {
                    [key: string]: components["schemas"]["VariableExpressionAO"];
                };
            };
        };
        responses: {
            /** @description Variables updated. New variables added and existing ones updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Environment with given `id` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    setEnvironmentVariables: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a environment
                 * @example 2v1av42-e525-4c00-a13a-1ac32d170724
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": {
                    [key: string]: components["schemas"]["VariableExpressionAO"];
                };
            };
        };
        responses: {
            /** @description Variables updated. New variables added, existing ones updated, all other removed. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Environment with given `id` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getExperiments: {
        parameters: {
            query?: {
                /**
                 * @description Filter results by experiments using the specified action. If multiple actions are specified, all of them needs to be used in the experiment
                 * @example com.steadybit.extension_host.stress-cpu
                 */
                action?: string[];
                /**
                 * @description Filter results by one or more external-ids
                 * @example incident-100
                 */
                externalId?: string[];
                /**
                 * @description Filter results via free text phrases searching for experiment name, key, property values, and 10 last run ids
                 * @example `ADM-1` or `#212`
                 */
                freeTextPhrases?: string[];
                /**
                 * @description Filter results by one or more experiments-keys
                 * @example ADM-1
                 */
                key?: string[];
                /**
                 * @description Filter results by experiments using an action with the specified kind. If multiple kinds are specified, all of them needs to be used in the experiment
                 * @example ATTACK
                 */
                kind?: ("ATTACK" | "CHECK" | "LOAD_TEST" | "OTHER" | "BASIC")[];
                /**
                 * @description Filter results by name and/or key of the experiment
                 * @example Outage
                 */
                name?: string;
                /**
                 * @description Filter results via properties
                 * @example `Value` or `EnumApiName:EnumValue`
                 */
                properties?: string[];
                /**
                 * @description Include only experiments which are runnable by the authorized user?
                 * @example false
                 */
                runnable?: boolean;
                /**
                 * @description Filter results by experiments linked to the given Service. If multiple services are specified, the experiment needs to be linked to all of them
                 * @example Shopping Cart Service
                 */
                service?: string[];
                /**
                 * @description Filter results by experiments having the specified tag. If multiple tags are specified, all of them needs to be assigned to the experiment
                 * @example kubernetes
                 */
                tag?: string[];
                /**
                 * @description Filter results by experiments using the specified target-type. If multiple target-types are specified, all of them needs to be used in the experiment
                 * @example com.steadybit.extension_host.host
                 */
                targetType?: string[];
                /**
                 * @description Filter results by one or more team-keys owning an experiment
                 * @example ADM
                 */
                team?: string[];
                /**
                 * @description Filter results by one or more team keys the experiment is shared with
                 * @example ADM
                 */
                teamSharedWith?: string[];
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains a summary of all existing experiments. Fetch a single experiment using the `experimentKey` to get more details of the experiment. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentSummariesAO"];
                    "application/yaml": components["schemas"]["ExperimentSummariesAO"];
                };
            };
        };
    };
    createOrUpdateExperiment: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CreateExperimentAO"];
            };
        };
        responses: {
            /** @description Experiment updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, experiment was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getExperiment: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
                 * @example ADM-2
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Experiment found, response body contains all the experiment details */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentAO"];
                    "application/yaml": components["schemas"]["ExperimentAO"];
                };
            };
            /** @description Experiment with given `key` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    updateExperiment: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
                 * @example ADM-2
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpdateExperimentAO"];
            };
        };
        responses: {
            /** @description Experiment updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, experiment was not updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteExperiment: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
                 * @example ADM-2
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Experiment deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment could not be found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getExperimentBadge: {
        parameters: {
            query: {
                /**
                 * @description Override the default hex-color `9f9f9f` for executions in state `canceled`.
                 * @example 9f9f9f
                 */
                colorMappingCanceled?: string;
                /**
                 * @description Override the default hex-color `4c1` for executions in state `completed`.
                 * @example 4c1
                 */
                colorMappingCompleted?: string;
                /**
                 * @description Override the default hex-color `fe7d37` for executions in state `created`.
                 * @example fe7d37
                 */
                colorMappingCreated?: string;
                /**
                 * @description Override the default hex-color `e05d44` for executions in state `errored`.
                 * @example e05d44
                 */
                colorMappingErrored?: string;
                /**
                 * @description Override the default hex-color `e05d44` for executions in state `failed`.
                 * @example e05d44
                 */
                colorMappingFailed?: string;
                /**
                 * @description Override the default hex-color `fe7d37` for executions in state `prepared`.
                 * @example fe7d37
                 */
                colorMappingPrepared?: string;
                /**
                 * @description Override the default hex-color `fe7d37` for executions in state `requested`.
                 * @example fe7d37
                 */
                colorMappingRequested?: string;
                /**
                 * @description Override the default hex-color `fe7d37` for executions in state `running`.
                 * @example fe7d37
                 */
                colorMappingRunning?: string;
                /**
                 * @description Optional parameter in case you need to scale the svg image. Defaults to 1
                 * @example 2
                 */
                scale?: number;
                /**
                 * @description Key of the Steadybit tenant (only for SaaS customers). You can get the key from the Platform URL or by asking the Steadybit team
                 * @example demo
                 */
                tenantKey: string;
            };
            header?: never;
            path: {
                /**
                 * @description Unique experiment key that identifies the experiment. Combination of `team key` and increasing number
                 * @example ADM-2
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Returns a badge as a svg-image */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "image/svg+xml": string;
                };
            };
        };
    };
    executeExperiment: {
        parameters: {
            query?: {
                /**
                 * @description By default an experiment is only executed when no other experiment is running. This can be overriden by starting the new experiment execution although another one is currently running
                 * @example true
                 */
                allowParallel?: boolean;
                /**
                 * @description Optional parameter to always store runs on any failure. If false, won´t be stored on validation errors (default behaviour).
                 * @example false
                 */
                forcePersist?: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique experiment key that identifies the experiment to be started. Combination of `team key` and increasing number
                 * @example ADM-2
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody?: {
            content: {
                "application/json": components["schemas"]["ExecuteExperimentRequestAO"];
            };
        };
        responses: {
            /** @description Experiment execution was started. The newly created object can be accessed via the HTTP header `location` */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExecuteExperimentResponseAO"];
                    "application/yaml": components["schemas"]["ExecuteExperimentResponseAO"];
                };
            };
            /** @description Experiment execution couldn't be started. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExecuteExperimentResponseAO"];
                    "application/yaml": components["schemas"]["ExecuteExperimentResponseAO"];
                };
            };
        };
    };
    getExperimentExecutions_3: {
        parameters: {
            query?: {
                /** @description Filter results by one or more states */
                state?: ("REQUESTED" | "CREATED" | "PREPARED" | "RUNNING" | "FAILED" | "CANCELED" | "COMPLETED" | "ERRORED")[];
            };
            header?: never;
            path: {
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains a summary of all experiment executions of the single experiment. Fetch a single experiment execution using the `id` of an execution  and the `/experiments/executions/{id}`-API to get more details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentExecutionSummariesAO"];
                    "application/yaml": components["schemas"]["ExperimentExecutionSummariesAO"];
                };
            };
        };
    };
    saveAndRun: {
        parameters: {
            query?: {
                /**
                 * @description Should this experiment also be executed when there is already at least one experiment running?
                 * @example false
                 */
                allowParallel?: boolean;
                /**
                 * @description Optional parameter to always store runs on any failure. If false, won´t be stored on validation errors (default behaviour).
                 * @example false
                 */
                forcePersist?: boolean;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CreateAndRunExperimentAO"];
            };
        };
        responses: {
            /** @description Experiment started */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExecuteExperimentResponseAO"];
                    "application/yaml": components["schemas"]["ExecuteExperimentResponseAO"];
                };
            };
            /** @description Validation error, experiment was neither saved nor executed */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExecuteExperimentResponseAO"];
                    "application/yaml": components["schemas"]["ExecuteExperimentResponseAO"];
                };
            };
        };
    };
    getExperimentExecutions_1: {
        parameters: {
            query?: {
                /**
                 * @description Filter results by name and/or key of the experiment
                 * @example Outage
                 */
                name?: string;
                /**
                 * @description Filter results by one or more states
                 * @example RUNNING
                 */
                state?: ("REQUESTED" | "CREATED" | "PREPARED" | "RUNNING" | "FAILED" | "CANCELED" | "COMPLETED" | "ERRORED")[];
                /**
                 * @description Filter results by one or more team-keys
                 * @example ADM
                 */
                team?: string[];
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains a summary of all experiment executions. Fetch a single experiment execution using the `id` of an execution  and the `/experiments/executions/{id}`-API to get more details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentExecutionSummariesAO"];
                    "application/yaml": components["schemas"]["ExperimentExecutionSummariesAO"];
                };
            };
        };
    };
    getExperimentExecutions_2: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ExperimentExecutionsRequestAO"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOExperimentExecutionPageItemAO"];
                    "application/yaml": components["schemas"]["PagedResponseAOExperimentExecutionPageItemAO"];
                };
            };
        };
    };
    getExperimentExecution: {
        parameters: {
            query?: {
                /**
                 * @description Additional fields to be returned for the experiment execution
                 * @example steps
                 */
                fields?: string;
            };
            header?: never;
            path: {
                /**
                 * @description Unique experiment execution id that identifies a single experiment execution
                 * @example 123
                 */
                id: number;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all information of the experiment execution of the single experiment. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentExecutionAO"];
                    "application/yaml": components["schemas"]["ExperimentExecutionAO"];
                };
            };
        };
    };
    getArtifact: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description The id of the artifact, usually the filename
                 * @example logfile.txt
                 */
                artifactId: string;
                /**
                 * @description Unique experiment execution id that identifies a single experiment execution
                 * @example 123
                 */
                id: number;
                /**
                 * @description The id of the target execution where the artifact is attached to
                 * @example 019abec0-3b14-7235-81fd-8007e720dfd0
                 */
                targetExecutionId: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    cancelExperimentExecution: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique experiment execution id that identifies a single experiment execution
                 * @example 123
                 */
                id: number;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Experiment execution was already canceled, errored or completed. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Request to cancel the experiment was accepted and will be performed by communicating to the experiment's agents. */
            202: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment execution not found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    updateExecutionProperties: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique experiment execution id that identifies a single experiment execution
                 * @example 123
                 */
                id: number;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpdateExperimentExecutionPropertiesAO"];
            };
        };
        responses: {
            /** @description Properties updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment execution not found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, properties were not updated. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    addExecutionPropertyValue: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique experiment execution id that identifies a single experiment execution
                 * @example 123
                 */
                id: number;
                /**
                 * @description The key of the property
                 * @example approvedBy
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": Record<string, never>;
            };
        };
        responses: {
            /** @description Properties updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment execution not found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, properties were not updated. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    setExecutionPropertyValue: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique experiment execution id that identifies a single experiment execution
                 * @example 123
                 */
                id: number;
                /**
                 * @description The key of the property
                 * @example approvedBy
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": Record<string, never>;
            };
        };
        responses: {
            /** @description Properties updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment execution not found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, properties were not updated. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    upsertSchedule: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertExperimentScheduleAO"];
            };
        };
        responses: {
            /** @description Schedule updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentScheduleAO"];
                    "application/yaml": components["schemas"]["ExperimentScheduleAO"];
                };
            };
            /** @description Schedule created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentScheduleAO"];
                    "application/yaml": components["schemas"]["ExperimentScheduleAO"];
                };
            };
            /** @description Validation error, schedule was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getSchedules: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Experiment schedule found. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentScheduleAO"];
                    "application/yaml": components["schemas"]["ExperimentScheduleAO"];
                };
            };
            /** @description Experiment Schedule couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    removeExperimentScheduleById: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Experiment schedule id.
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Experiment schedule deleted. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment Schedule couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    patchSchedule: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Experiment schedule id.
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PatchExperimentScheduleAO"];
            };
        };
        responses: {
            /** @description Schedule updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentScheduleAO"];
                    "application/yaml": components["schemas"]["ExperimentScheduleAO"];
                };
            };
            /** @description Experiment Schedule couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, schedule was not updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getAllSchedulesV2: {
        parameters: {
            query?: {
                /**
                 * @description Filter results by one or more experiment-keys
                 * @example ADM-5
                 */
                experiment?: string[];
                /**
                 * @description Filter results by one or more team-keys
                 * @example ADM
                 */
                team?: string[];
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description All experiment schedule configurations */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentScheduleAO"][];
                    "application/yaml": components["schemas"]["ExperimentScheduleAO"][];
                };
            };
        };
    };
    getExperimentTemplates: {
        parameters: {
            query?: {
                /** @description Filter results by one or more action, like `com.steadybit.extension_host.stress-cpu` */
                action?: string[];
                /** @description Filter results by one or more free text phrases searching in the template title and template description */
                freeTextPhrases?: string[];
                /**
                 * @description Include hidden templates (requires an admin token)
                 * @example false
                 */
                includeHidden?: boolean;
                /**
                 * @description Include templates referencing actions/target-types/property-definitions that are not available
                 * @example false
                 */
                includeNonAvailable?: boolean;
                /** @description Filter results by one or more tags */
                tag?: string[];
                /** @description Filter results by one or more target type, like `com.steadybit.extension_container.container` */
                targetType?: string[];
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all existing templates. Fetch a single template by `id` to get more details for a template */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentTemplateSummariesAO"];
                    "application/yaml": components["schemas"]["ExperimentTemplateSummariesAO"];
                };
            };
        };
    };
    upsertExperimentTemplate: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertExperimentTemplateAO"];
            };
        };
        responses: {
            /** @description Experiment Template updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentTemplateAO"];
                    "application/yaml": components["schemas"]["ExperimentTemplateAO"];
                };
            };
            /** @description Experiment Template created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentTemplateAO"];
                    "application/yaml": components["schemas"]["ExperimentTemplateAO"];
                };
            };
            /** @description Validation error, experiment template was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getExperimentTemplate: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of an experiment template
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Experiment Template found, response body contains the team details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExperimentTemplateAO"];
                    "application/yaml": components["schemas"]["ExperimentTemplateAO"];
                };
            };
            /** @description Experiment Template with given `id` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteExperimentTemplate: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of an experiment template
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Experiment Template deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment Template with given `id` was not found and thus not deleted */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    createExperimentByTemplate: {
        parameters: {
            query?: {
                /**
                 * @description If `true`, all properties will be reset to properties specified in the template either with their fixed values in the template or via template placeholder. If `false`, existing properties will stay untouched, only new properties will be added. Only relevant for experiment updates via `externalId`.
                 * @example true
                 */
                resetProperties?: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique identifier of an experiment template
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CreateExperimentFromTemplateAO"];
            };
        };
        responses: {
            /** @description Experiment updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Template specified by experimentTemplateId could not be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, experiment was not created */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    saveAndRunFromTemplate: {
        parameters: {
            query?: {
                /**
                 * @description Should this experiment also be executed when there is already another experiment running?
                 * @example false
                 */
                allowParallel?: boolean;
                /**
                 * @description Optional parameter to always store runs on any failure. If false, won´t be stored on validation errors (default behaviour).
                 * @example false
                 */
                forcePersist?: boolean;
                /**
                 * @description If `true`, all properties will be reset to properties specified in the template either with their fixed values in the template or via template placeholder. If `false`, existing properties will stay untouched, only new properties will be added. Only relevant for experiment updates via `externalId`.
                 * @example true
                 */
                resetProperties?: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique identifier of an experiment template
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CreateAndRunExperimentFromTemplateAO"];
            };
        };
        responses: {
            /** @description Experiment started */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExecuteExperimentResponseAO"];
                    "application/yaml": components["schemas"]["ExecuteExperimentResponseAO"];
                };
            };
            /** @description Validation error, experiment was neither saved nor executed */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ExecuteExperimentResponseAO"];
                    "application/yaml": components["schemas"]["ExecuteExperimentResponseAO"];
                };
            };
        };
    };
    updateExperimentByTemplate: {
        parameters: {
            query?: {
                /**
                 * @description If `true`, all properties will be reset to properties specified in the template either with their fixed values in the template or via template placeholder. If `false`, existing properties will stay untouched, only new properties will be added.
                 * @example true
                 */
                resetProperties?: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique identifier of an experiment template
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
                /**
                 * @description The key of the experiment that should be updated
                 * @example ADM-18
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpdateExperimentFromTemplateAO"];
            };
        };
        responses: {
            /** @description Experiment updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Template specified by experimentTemplateId or experiment specified by key could not be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, experiment was not created */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    importFromHub: {
        parameters: {
            query?: {
                /**
                 * @description Do you want to overwrite a template that already exists? If set to `false` and any of the templates already exist, the API will return HTTP status 409 and none of the templates are imported.
                 * @example true
                 */
                overwrite?: boolean;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ExperimentTemplatesImportAO"];
            };
        };
        responses: {
            /** @description Experiment Templates imported */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description If any of the templates to be imported already exist in the platform and overwrite is set to `false` */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getLandscapeViews: {
        parameters: {
            query: {
                /**
                 * @description Key of the team whose saved views should be returned.
                 * @example ADM
                 */
                team: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all saved landscape views of the team. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListResponseLandscapeViewAO"];
                    "application/yaml": components["schemas"]["ListResponseLandscapeViewAO"];
                };
            };
            /** @description The given team could not be found. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    createLandscapeView: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertLandscapeViewAO"];
            };
        };
        responses: {
            /** @description Saved view created, response body contains its details. */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["LandscapeViewAO"];
                    "application/yaml": components["schemas"]["LandscapeViewAO"];
                };
            };
            /** @description Validation error, the saved view was not created. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getLandscapeView: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a saved view.
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Saved view found, response body contains its details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["LandscapeViewAO"];
                    "application/yaml": components["schemas"]["LandscapeViewAO"];
                };
            };
            /** @description The given `id` is not a uuid. */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Saved view with given `id` was not found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    updateLandscapeView: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a saved view.
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertLandscapeViewAO"];
            };
        };
        responses: {
            /** @description Saved view updated, response body contains its details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["LandscapeViewAO"];
                    "application/yaml": components["schemas"]["LandscapeViewAO"];
                };
            };
            /** @description Saved view with given `id` was not found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, the saved view was not updated. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteLandscapeView: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a saved view.
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Saved view deleted. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Saved view with given `id` was not found and thus not deleted. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    health: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": Record<string, never>;
                    "application/yaml": Record<string, never>;
                };
            };
        };
    };
    liveness: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": Record<string, never>;
                    "application/yaml": Record<string, never>;
                };
            };
        };
    };
    readiness: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": Record<string, never>;
                    "application/yaml": Record<string, never>;
                };
            };
        };
    };
    getHubs: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all hubs. Fetch a single hub by `id` to get more details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["HubSummariesAO"];
                    "application/yaml": components["schemas"]["HubSummariesAO"];
                };
            };
        };
    };
    upsertHub: {
        parameters: {
            query?: {
                /** @description Whether to synchronize the hub or not. */
                synchronize?: boolean;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertHubAO"];
            };
        };
        responses: {
            /** @description Hub updated. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["HubAO"];
                    "application/yaml": components["schemas"]["HubAO"];
                };
            };
            /** @description Hub created. */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["HubAO"];
                    "application/yaml": components["schemas"]["HubAO"];
                };
            };
            /** @description Validation error, hub was not created / updated. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getHubById: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a hub
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Hub found, response body contains the hub. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["HubAO"];
                    "application/yaml": components["schemas"]["HubAO"];
                };
            };
            /** @description Hub with the given `id` was not found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteHub: {
        parameters: {
            query?: {
                /** @description Whether imported templates of the hub should be deleted as well. */
                deleteImportedTemplates?: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a hub.
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Hub deleted. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Hub with given `id` was not found and thus not deleted. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    resyncHub: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a hub
                 * @example d7e65100-1d20-4980-be87-c351704910b8
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Hub re-synchronized. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["HubAO"];
                    "application/yaml": components["schemas"]["HubAO"];
                };
            };
            /** @description Hub with given `id` was not found and thus not re-synchronized. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Rate limit exceeded. */
            429: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    connectionCheck: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["HubConnectionCheckAO"];
            };
        };
        responses: {
            /** @description Hub successful resolved. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["HubConnectionCheckResponseAO"];
                    "application/yaml": components["schemas"]["HubConnectionCheckResponseAO"];
                };
            };
            /** @description Hub could not be resolved. */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getPreflightWebhooks: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all existing preflight webhooks. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListResponsePreflightWebhookAO"];
                    "application/yaml": components["schemas"]["ListResponsePreflightWebhookAO"];
                };
            };
        };
    };
    upsertPreflightWebhook: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PreflightWebhookUpsertAO"];
            };
        };
        responses: {
            /** @description Preflight webhook updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PreflightWebhookAO"];
                    "application/yaml": components["schemas"]["PreflightWebhookAO"];
                };
            };
            /** @description Preflight webhook created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PreflightWebhookAO"];
                    "application/yaml": components["schemas"]["PreflightWebhookAO"];
                };
            };
            /** @description Validation error, preflight webhook was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getPreflightActionIntegrations: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all existing preflight preflight action integrations. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListResponsePreflightActionIntegrationAO"];
                    "application/yaml": components["schemas"]["ListResponsePreflightActionIntegrationAO"];
                };
            };
        };
    };
    upsertPreflightActionIntegration: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PreflightActionIntegrationUpsertAO"];
            };
        };
        responses: {
            /** @description Preflight action integration updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PreflightActionIntegrationAO"];
                    "application/yaml": components["schemas"]["PreflightActionIntegrationAO"];
                };
            };
            /** @description Preflight action integration created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PreflightActionIntegrationAO"];
                    "application/yaml": components["schemas"]["PreflightActionIntegrationAO"];
                };
            };
            /** @description Validation error, preflight action integration was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getPreflightActionIntegration: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description ID of the preflight action integration
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Preflight preflight action integrations found, response body contains the details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PreflightActionIntegrationAO"];
                    "application/yaml": components["schemas"]["PreflightActionIntegrationAO"];
                };
            };
            /** @description Preflight preflight action integrations with given `id` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deletePreflightActionIntegration: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a preflight action integration
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Preflight action integration deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PreflightActionIntegrationAO"];
                    "application/yaml": components["schemas"]["PreflightActionIntegrationAO"];
                };
            };
            /** @description Preflight action integration with given `id` was not found and thus not deleted */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getPreflightWebhook: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description ID of the preflight webhook
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Preflight webhooks found, response body contains the details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PreflightWebhookAO"];
                    "application/yaml": components["schemas"]["PreflightWebhookAO"];
                };
            };
            /** @description Preflight webhooks with given `id` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deletePreflightWebhook: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a preflight webhook
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Preflight webhook deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PreflightWebhookAO"];
                    "application/yaml": components["schemas"]["PreflightWebhookAO"];
                };
            };
            /** @description Preflight webhook with given `id` was not found and thus not deleted */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getSlackIntegrations: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all existing integrations. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListResponseSlackWebhookAO"];
                    "application/yaml": components["schemas"]["ListResponseSlackWebhookAO"];
                };
            };
        };
    };
    upsertSlackIntegration: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["SlackWebhookUpsertAO"];
            };
        };
        responses: {
            /** @description Slack integration updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SlackWebhookAO"];
                    "application/yaml": components["schemas"]["SlackWebhookAO"];
                };
            };
            /** @description Slack integration created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SlackWebhookAO"];
                    "application/yaml": components["schemas"]["SlackWebhookAO"];
                };
            };
            /** @description Validation error, slack integration was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getSlackIntegration: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description ID of the integration
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Integration found, response body contains the details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SlackWebhookAO"];
                    "application/yaml": components["schemas"]["SlackWebhookAO"];
                };
            };
            /** @description Integration with given `id` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteSlackIntegration: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a slack integration
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Slack integration deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SlackWebhookAO"];
                    "application/yaml": components["schemas"]["SlackWebhookAO"];
                };
            };
            /** @description Slack integration with given `id` was not found and thus not deleted */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getCustomWebhooks: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all existing custom webhooks. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListResponseCustomWebhookAO"];
                    "application/yaml": components["schemas"]["ListResponseCustomWebhookAO"];
                };
            };
        };
    };
    upsertCustomWebhook: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["CustomWebhookUpsertAO"];
            };
        };
        responses: {
            /** @description Custom webhook updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CustomWebhookAO"];
                    "application/yaml": components["schemas"]["CustomWebhookAO"];
                };
            };
            /** @description Custom webhook created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CustomWebhookAO"];
                    "application/yaml": components["schemas"]["CustomWebhookAO"];
                };
            };
            /** @description Validation error, custom webhook was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getCustomWebhook: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description ID of the custom webhook
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Custom webhooks found, response body contains the details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CustomWebhookAO"];
                    "application/yaml": components["schemas"]["CustomWebhookAO"];
                };
            };
            /** @description Custom webhooks with given `id` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteCustomWebhook: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a custom webhook
                 * @example ac456d58-8fb2-4df4-86d8-ca81d7562739
                 */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Custom webhook deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CustomWebhookAO"];
                    "application/yaml": components["schemas"]["CustomWebhookAO"];
                };
            };
            /** @description Custom webhook with given `id` was not found and thus not deleted */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getKillswitch: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Kill switch status in the response body */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["KillswitchAO"];
                    "application/yaml": components["schemas"]["KillswitchAO"];
                };
            };
        };
    };
    engageKillswitch: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Kill switch was activated / engaged */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    disengageKillswitch: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Kill switch was deactivated / disengaged */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getLicenseSummary: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description License summary. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["GetLicenseSummaryAO"];
                    "application/yaml": components["schemas"]["GetLicenseSummaryAO"];
                };
            };
        };
    };
    getReport: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description License report zip. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getPreflightActionSummary: {
        parameters: {
            query: {
                offset: number;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Matching targets with preflightAction information, limited to 20 items. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PreflightActionSummaryAO"];
                    "application/yaml": components["schemas"]["PreflightActionSummaryAO"];
                };
            };
        };
    };
    getAssociations: {
        parameters: {
            query?: {
                /**
                 * @description Filter association based on association type (`EXPERIMENT` for experiment-related associations and `SERVICE` for service-related associations, no matter whether globally or individually)
                 * @example EXPERIMENT
                 */
                associationTypeAO?: "EXPERIMENT" | "SERVICE";
                /**
                 * @description Filter association that are explicitly assigned to the given experimentKey. (There might still be associations for ALL Experiment Designs)
                 * @example ADM-15
                 */
                experimentKey?: string;
                /**
                 * @description Filter association based on a single property definition key
                 * @example RESULT_COLOR
                 */
                key?: string;
                /**
                 * @description The number of the page, responses are limited to 50 elements per page.
                 * @example 0
                 */
                page?: number;
                /**
                 * @description Filter association that are explicitly assigned to the given serviceId. (There might still be associations for ALL Services)
                 * @example c1975ae7-02f6-4a4a-9762-14559b330b8c
                 */
                serviceId?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description List of associations. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOPropertyAssociationAO"];
                    "application/yaml": components["schemas"]["PagedResponseAOPropertyAssociationAO"];
                };
            };
        };
    };
    upsertPropertyAssociation: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertPropertyAssociationAO"];
            };
        };
        responses: {
            /** @description Property association updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PropertyAssociationAO"];
                    "application/yaml": components["schemas"]["PropertyAssociationAO"];
                };
            };
            /** @description Property association created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PropertyAssociationAO"];
                    "application/yaml": components["schemas"]["PropertyAssociationAO"];
                };
            };
            /** @description Version not matching */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, property association was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getPropertyDefinition_1: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Property association found. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PropertyAssociationAO"];
                    "application/yaml": components["schemas"]["PropertyAssociationAO"];
                };
            };
            /** @description Property association couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deletePropertyAssociation: {
        parameters: {
            query?: {
                /**
                 * @description Associations can only be deleted, if no experiment design or experiment schedule is still using the value. Setting this parameter to `true` will delete those values.
                 * @example false
                 */
                deleteValues?: boolean;
            };
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Property association deleted. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Property association couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getPropertyDefinitions: {
        parameters: {
            query: {
                /**
                 * @description The number of the page, responses are limited to 50 elements per page.
                 * @example 0
                 */
                page: components["schemas"]["PageRequestAO"];
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOPropertyDefinitionAO"];
                    "application/yaml": components["schemas"]["PagedResponseAOPropertyDefinitionAO"];
                };
            };
        };
    };
    upsertPropertyDefinition: {
        parameters: {
            query?: {
                /**
                 * @description You can remove enum-values for a ENUM or ENUM_LIST property if they are still in use in experiment designs. Setting this parameter to `true` will delete those values in experiment designs.
                 * @example false
                 */
                deleteValues?: boolean;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertPropertyDefinitionAO"];
            };
        };
        responses: {
            /** @description Property definition updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PropertyDefinitionAO"];
                    "application/yaml": components["schemas"]["PropertyDefinitionAO"];
                };
            };
            /** @description Property definition created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PropertyDefinitionAO"];
                    "application/yaml": components["schemas"]["PropertyDefinitionAO"];
                };
            };
            /** @description Version not matching */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, property definition was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getPropertyDefinition: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Property definition found. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PropertyDefinitionAO"];
                    "application/yaml": components["schemas"]["PropertyDefinitionAO"];
                };
            };
            /** @description Property definition couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deletePropertyDefinition: {
        parameters: {
            query?: {
                /**
                 * @description Definitions can only be deleted, if no associations are still refering to this property. Setting the value to `true` will delete all associations and all current values in experiment designs and schedules. Existing executions won't get touched.
                 * @example false
                 */
                deleteAssociations?: boolean;
            };
            header?: never;
            path: {
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Property definition deleted. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Property definition couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getEnvironmentCounts: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ReportFilterAO"];
            };
        };
        responses: {
            /** @description Environment count time series. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TimeSeriesReportAO"];
                    "application/yaml": components["schemas"]["TimeSeriesReportAO"];
                };
            };
        };
    };
    getExperimentCreations: {
        parameters: {
            query?: {
                /**
                 * @description Grouping dimension for the results.
                 * @example CREATED_VIA
                 */
                groupBy?: "NONE" | "CREATED_VIA" | "ORIGIN";
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ExperimentReportFilterAO"];
            };
        };
        responses: {
            /** @description Creation count time series. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TimeSeriesReportAO"];
                    "application/yaml": components["schemas"]["TimeSeriesReportAO"];
                };
            };
        };
    };
    getExperimentExecutions: {
        parameters: {
            query?: {
                /**
                 * @description Grouping dimension for the results.
                 * @example STATE
                 */
                groupBy?: "NONE" | "STATE" | "TRIGGER" | "ACTION" | "ISSUES_FIXED" | "ISSUES_DISCOVERED";
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ExperimentExecutionReportFilterAO"];
            };
        };
        responses: {
            /** @description Execution count time series. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TimeSeriesReportAO"];
                    "application/yaml": components["schemas"]["TimeSeriesReportAO"];
                };
            };
        };
    };
    getAverageRisk: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ServiceRiskReportFilterAO"];
            };
        };
        responses: {
            /** @description Average risk time series. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TimeSeriesReportAO"];
                    "application/yaml": components["schemas"]["TimeSeriesReportAO"];
                };
            };
        };
    };
    getRiskByCategory: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ServiceRiskReportFilterAO"];
            };
        };
        responses: {
            /** @description Risk by category time series. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TimeSeriesReportAO"];
                    "application/yaml": components["schemas"]["TimeSeriesReportAO"];
                };
            };
        };
    };
    getRiskDistribution: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ServiceRiskReportFilterAO"];
            };
        };
        responses: {
            /** @description Risk level distribution time series. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TimeSeriesReportAO"];
                    "application/yaml": components["schemas"]["TimeSeriesReportAO"];
                };
            };
        };
    };
    getTeamCounts: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ReportFilterAO"];
            };
        };
        responses: {
            /** @description Team count time series. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TimeSeriesReportAO"];
                    "application/yaml": components["schemas"]["TimeSeriesReportAO"];
                };
            };
        };
    };
    getUserCounts: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["ReportFilterAO"];
            };
        };
        responses: {
            /** @description User count time series. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TimeSeriesReportAO"];
                    "application/yaml": components["schemas"]["TimeSeriesReportAO"];
                };
            };
        };
    };
    getServiceList: {
        parameters: {
            query: {
                /** @description Filter results by one or more environment name, like 'Global' */
                environmentName?: string[];
                /** @description Filter results by one or more experiment keys being linked to a service, like 'ADM-123' */
                experimentKey?: string[];
                page: components["schemas"]["PageRequestAO"];
                /** @description Filter results by one or more team key, like 'ADM' */
                teamKey?: string[];
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body containing all existing services matching the query parameters. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOServiceSummaryAO"];
                    "application/yaml": components["schemas"]["PagedResponseAOServiceSummaryAO"];
                };
            };
        };
    };
    upsertService: {
        parameters: {
            query?: {
                /**
                 * @description When the service profile of a service gets updated, provided experiments whose template ids are not part of the new service profile are not allowed. Setting the value to `true` will delete those experiments.
                 * @example false
                 */
                deleteExperiments?: boolean;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertServiceAO"];
            };
        };
        responses: {
            /** @description Service updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ServiceAO"];
                    "application/yaml": components["schemas"]["ServiceAO"];
                };
            };
            /** @description Service created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ServiceAO"];
                    "application/yaml": components["schemas"]["ServiceAO"];
                };
            };
            /** @description Version not matching */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, Service was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getService: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Service found. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ServiceAO"];
                    "application/yaml": components["schemas"]["ServiceAO"];
                };
            };
            /** @description Service couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteService: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Service deleted. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Service couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getServiceExperiments: {
        parameters: {
            query: {
                /** @description Filter results by one or more categories */
                category?: string[];
                /** @description Include custom experiments with missing categories */
                categoryMissing?: boolean;
                page: components["schemas"]["PageRequestAO"];
                /** @description Filter results by type (PROVIDED,CUSTOM) */
                type?: ("PROVIDED" | "CUSTOM")[];
            };
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Service-Associations found. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOServiceExperimentAO"];
                    "application/yaml": components["schemas"]["PagedResponseAOServiceExperimentAO"];
                };
            };
        };
    };
    linkCustomExperiment: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["LinkCustomExperimentRequestAO"];
            };
        };
        responses: {
            /** @description Experiment linked to service. */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    unlinkCustomExperiment: {
        parameters: {
            query: {
                experimentKey: string;
            };
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Experiment unlinked from service. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    upsertProvidedExperiment: {
        parameters: {
            query?: {
                /**
                 * @description If `true`, all properties will be reset to properties specified in the template either with their fixed values in the template or via template placeholder. If `false`, existing properties will stay untouched, only new properties will be added. Only relevant for experiment updates.
                 * @example true
                 */
                resetProperties?: boolean;
            };
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertProvidedExperimentRequestAO"];
            };
        };
        responses: {
            /** @description Experiment updated. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Experiment created. */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Service or Template not found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getRisk: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Risk score found. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ServiceRiskAO"];
                    "application/yaml": components["schemas"]["ServiceRiskAO"];
                };
            };
            /** @description Service or risk not found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getServiceVariables: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Unique identifier of a service */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Key-value map of variables owned by the service. Each value is either a constant string, an array of constant strings, or a select expression object. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": string;
                    "application/yaml": string;
                };
            };
            /** @description Service not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    setServiceVariables: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Unique identifier of a service */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": {
                    [key: string]: components["schemas"]["VariableExpressionAO"];
                };
            };
        };
        responses: {
            /** @description Variables updated. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Service not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    mergeServiceVariables: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Unique identifier of a service */
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": {
                    [key: string]: components["schemas"]["VariableExpressionAO"];
                };
            };
        };
        responses: {
            /** @description Variables updated. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Service not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getProfiles: {
        parameters: {
            query: {
                /** @description Filter results by defaultProfile flag */
                defaultProfile?: boolean;
                /** @description Filter results by name (partial match) */
                name?: string;
                /** @description Filter results by origin (PROVIDED, CUSTOM) */
                origin?: string[];
                page: components["schemas"]["PageRequestAO"];
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body containing all existing service profiles matching the query parameters. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOServiceProfileAO"];
                    "application/yaml": components["schemas"]["PagedResponseAOServiceProfileAO"];
                };
            };
        };
    };
    upsertProfile: {
        parameters: {
            query?: {
                /**
                 * @description When templates are removed from a service profile, provided experiments using those templates will be affected. Setting the value to `true` will delete those experiments.
                 * @example false
                 */
                deleteExperiments?: boolean;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertServiceProfileAO"];
            };
        };
        responses: {
            /** @description Service profile updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ServiceProfileAO"];
                    "application/yaml": components["schemas"]["ServiceProfileAO"];
                };
            };
            /** @description Service profile created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ServiceProfileAO"];
                    "application/yaml": components["schemas"]["ServiceProfileAO"];
                };
            };
            /** @description Version not matching */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, Service profile was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getProfile: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Service profile found. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ServiceProfileAO"];
                    "application/yaml": components["schemas"]["ServiceProfileAO"];
                };
            };
            /** @description Service profile couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteProfile: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Service profile deleted. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Service profile couldn't be found. */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Cannot delete PROVIDED profiles. */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getTargetsStats: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": {
                        [key: string]: number;
                    };
                    "application/yaml": {
                        [key: string]: number;
                    };
                };
            };
        };
    };
    getTargetsStats_1: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["TargetStatsRequest"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": {
                        [key: string]: number;
                    };
                    "application/yaml": {
                        [key: string]: number;
                    };
                };
            };
        };
    };
    getTargets: {
        parameters: {
            query: {
                /** @description Optional, list of requested target attribute keys. If not specified, all attributes will be returned. Multiple values allowed. Example: `k8s.deployment` */
                attribute?: string[];
                /**
                 * @description Optional, the cursor to use to fetch the next page
                 * @example eyJhZ2VudElkIjoiMDE5ZDE4ZjYtZTgwYy03NDdlLThkYWItNmE5MTBkM2JhZWQyIiwibmFtZSI6ImRlbW8tZGV2ZWxvcC9pbmdyZXNzLW5naW54L2luZ3Jlc3MtbmdpbngtY29udHJvbGxlciIsInR5cGUiOiJjb20uc3RlYWR5Yml0LmV4dGVuc2lvbl9rdWJlcm5ldGVzLmt1YmVybmV0ZXMtZGVwbG95bWVudCJ9
                 */
                cursor?: string;
                /**
                 * @description The name of the environment
                 * @example Global
                 */
                environment: string;
                /**
                 * @description Optional, additional target selection query
                 * @example (k8s.cluster-name="demo-develop" AND k8s.namespace="steadybit-demo" AND k8s.deployment="fashion-bestseller")
                 */
                query?: string;
                /** @description Optional, the number of items to return per page. default is 100, maximum is 1000. */
                size?: number;
                /**
                 * @description Optional, the type of the target
                 * @example com.steadybit.extension_kubernetes.kubernetes-deployment
                 */
                targetType?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description List of targets. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["CursorSliceResponseAOTargetAO"];
                    "application/yaml": components["schemas"]["CursorSliceResponseAOTargetAO"];
                };
            };
            /** @description Environment not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getTargetAttributeKeys: {
        parameters: {
            query: {
                /**
                 * @description If the action specifies a extended target selector and you want to fetch all attribute keys for a given action. Required if targetType is not set
                 * @example com.steadybit.extension_kubernetes.delete-pod
                 */
                actionId?: string;
                /**
                 * @description The name of the environment
                 * @example Global
                 */
                environment: string;
                /** @description The page number to retrieve. Starts from 0. default is 0 */
                page?: number;
                /** @description The number of items to return per page. default is 100, maximum is 100. */
                size?: number;
                /**
                 * @description The type of the target, required if actionId is not set
                 * @example com.steadybit.extension_kubernetes.kubernetes-deployment
                 */
                targetType?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description List of attribute keys. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOString"];
                    "application/yaml": components["schemas"]["PagedResponseAOString"];
                };
            };
            /** @description Environment not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getTargetAttributeValues: {
        parameters: {
            query: {
                /**
                 * @description If the action specifies a extended target selector and you want to fetch all attribute keys for a given action. Required if targetType is not set
                 * @example com.steadybit.extension_kubernetes.delete-pod
                 */
                actionId?: string;
                /**
                 * @description The key of of the attribute
                 * @example k8s.namespace
                 */
                attributeKey: string;
                /**
                 * @description The name of the environment
                 * @example Global
                 */
                environment: string;
                /** @description The page number to retrieve. Starts from 0. default is 0 */
                page?: number;
                /** @description The number of items to return per page. default is 100, maximum is 100. */
                size?: number;
                /**
                 * @description The type of the target, required if actionId is not set
                 * @example com.steadybit.extension_kubernetes.kubernetes-deployment
                 */
                targetType?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description List of attribute values. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["PagedResponseAOString"];
                    "application/yaml": components["schemas"]["PagedResponseAOString"];
                };
            };
            /** @description Environment not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getTeams: {
        parameters: {
            query?: {
                /**
                 * @description If set and used with an `accessToken` associated to one or multiple teams, only the team associated to the token are returned. Otherwise, all teams are listed.
                 * @example true
                 */
                onlyAccessible?: boolean;
                /**
                 * @description If set, only teams matching the search are returned. Matches the team name or key, the name or email of a team member, or the name of an allowed environment.
                 * @example shop
                 */
                search?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Response body contains all existing teams. Fetch a single team by `key` to get more details for a team */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamSummariesAO"];
                    "application/yaml": components["schemas"]["TeamSummariesAO"];
                };
            };
        };
    };
    upsertTeam: {
        parameters: {
            query?: {
                /**
                 * @description By default, Steadybit checks whether the allowed actions exists and are reported by an agent. For convenience, this can be deactivated to decouple team creation and agent-installation
                 * @example false
                 */
                validateActions?: boolean;
                /**
                 * @description By default, Steadybit will skip members, which are not yet know. If set to true, Steadybit will validate the given members and show a 422 response.
                 * @example false
                 */
                validateMembers?: boolean;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["UpsertTeamAO"];
            };
        };
        responses: {
            /** @description Team updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamAO"];
                    "application/yaml": components["schemas"]["TeamAO"];
                };
            };
            /** @description Team created */
            201: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamAO"];
                    "application/yaml": components["schemas"]["TeamAO"];
                };
            };
            /** @description Validation error, team was not created / updated */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getTeam: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Team found, response body contains the team details. */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamAO"];
                    "application/yaml": components["schemas"]["TeamAO"];
                };
            };
            /** @description Team with given `key` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    deleteTeam: {
        parameters: {
            query: {
                /**
                 * @description Safety-Parameter - purge team including all experiments and executions.
                 * @example true
                 */
                purgeIncludingExperiments: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Team deleted */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamAO"];
                    "application/yaml": components["schemas"]["TeamAO"];
                };
            };
            /** @description Team has running experiments or insufficient permisssions */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamAO"];
                    "application/yaml": components["schemas"]["TeamAO"];
                };
            };
            /** @description Team with given `key` was not found and thus not deleted */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getTeamEnvironments: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Team found, response body contains all assigned environments */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamEnvironmentsAO"];
                    "application/yaml": components["schemas"]["TeamEnvironmentsAO"];
                };
            };
            /** @description Team with given `key` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    setTeamEnvironments: {
        parameters: {
            query?: {
                /**
                 * @description By default, Steadybit will skip environments, which are not yet know. If set to true, Steadybit will validate the given environments and show a 422 response.
                 * @example false
                 */
                validateEnvironments?: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        /** @description Update request to change the environments of a specific team. */
        requestBody: {
            content: {
                "application/json": components["schemas"]["TeamEnvironmentsAO"];
            };
        };
        responses: {
            /** @description Team environments updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamEnvironmentsAO"];
                    "application/yaml": components["schemas"]["TeamEnvironmentsAO"];
                };
            };
            /** @description Validation error, no update was performed */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    addTeamEnvironments: {
        parameters: {
            query?: {
                /**
                 * @description By default, Steadybit will skip environments, which are not yet know. If set to true, Steadybit will validate the given environments and show a 422 response.
                 * @example false
                 */
                validateEnvironments?: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["TeamEnvironmentsUpdateAO"];
            };
        };
        responses: {
            /** @description Team environments added */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamEnvironmentsAO"];
                    "application/yaml": components["schemas"]["TeamEnvironmentsAO"];
                };
            };
            /** @description Validation error, no team environment were added */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    removeTeamEnvironments: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["TeamEnvironmentsUpdateAO"];
            };
        };
        responses: {
            /** @description Team environments removed */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamEnvironmentsAO"];
                    "application/yaml": components["schemas"]["TeamEnvironmentsAO"];
                };
            };
            /** @description Validation error, no team environments were removed */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    getTeamMembers: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Team found, response body contains all the team members */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamMembersAO"];
                    "application/yaml": components["schemas"]["TeamMembersAO"];
                };
            };
            /** @description Team with given `key` was not found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    setTeamMembers: {
        parameters: {
            query?: {
                /**
                 * @description By default, Steadybit will skip members, which are not yet know. If set to true, Steadybit will validate the given members and show a 422 response.
                 * @example false
                 */
                validateMembers?: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        /** @description Update request to change the members of a specific team. Specify either username, being a Steadybit user id, or the email address of the user. */
        requestBody: {
            content: {
                "application/json": components["schemas"]["TeamMembersUpdateAO"];
            };
        };
        responses: {
            /** @description Team members updated */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamMembersAO"];
                    "application/yaml": components["schemas"]["TeamMembersAO"];
                };
            };
            /** @description Validation error, no update was performed */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    addTeamMembers: {
        parameters: {
            query?: {
                /**
                 * @description By default, Steadybit will skip members, which are not yet know. If set to true, Steadybit will validate the given members and show a 422 response.
                 * @example false
                 */
                validateMembers?: boolean;
            };
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["TeamMembersUpdateAO"];
            };
        };
        responses: {
            /** @description Team members added */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamMembersAO"];
                    "application/yaml": components["schemas"]["TeamMembersAO"];
                };
            };
            /** @description Validation error, no team member were added */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    removeTeamMembers: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /**
                 * @description Unique identifier of a team
                 * @example ADM
                 */
                key: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["TeamMembersRemoveAO"];
            };
        };
        responses: {
            /** @description Team members removed */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TeamMembersAO"];
                    "application/yaml": components["schemas"]["TeamMembersAO"];
                };
            };
            /** @description Validation error, no team members were removed */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
    inviteUser: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["InviteUsersRequestAO"];
            };
        };
        responses: {
            /** @description Users invited */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Insufficient permissions to invite a new user to this tenant. */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Validation error, no or not all users were invited */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
        };
    };
}
