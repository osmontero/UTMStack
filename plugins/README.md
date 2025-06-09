# ThreatWinds EventProcessor and UTMStack Integration

This documentation provides a comprehensive guide on how to implement rules for analysis, and pipelines for data
extraction, enrichment, and transformation within the EventProcessor and UTMStack ecosystem. It is designed to be a
practical reference for developers working with these systems.

## Table of Contents

1. [Introduction](#introduction)
2. [Architecture Overview](#architecture-overview)
3. [Components](#components)
4. [Implementing Rules](#implementing-rules)
5. [Implementing Filters](#implementing-filters)
6. [Best Practices](#best-practices)
7. [Development Workflow](#development-workflow)
8. [Advanced Features](#advanced-features)
9. [Integration with Other Systems](#integration-with-other-systems)
10. [Performance Optimization](#performance-optimization)
11. [Troubleshooting](#troubleshooting)
12. [Real-World Use Cases](#real-world-use-cases)
13. [Custom Plugin Development](#custom-plugin-development)
14. [Scaling the System](#scaling-the-system)
15. [Migration and Upgrades](#migration-and-upgrades)
16. [Community Resources and Support](#community-resources-and-support)

## Introduction

This documentation provides a comprehensive guide on how to implement rules for analysis, extraction, and data
transformation within the EventProcessor and UTMStack ecosystem. The EventProcessor is a security event processing
engine that uses a plugin-based architecture to process, analyze, and transform security events from various sources.
UTMStack integrates with the EventProcessor to provide a complete security monitoring solution.

## Architecture Overview

The EventProcessor and UTMStack integration consists of several components:

1. **EventProcessor**: The core engine that processes security events.
2. **EventProcessor Plugins**: Plugins that extend the functionality of the EventProcessor.
3. **UTMStack Plugins**: Plugins that integrate UTMStack with the EventProcessor.
4. **go-sdk**: A Go SDK that provides common functionality for both the EventProcessor and plugins.
5. **Rules**: YAML files that define analysis rules for detecting security threats.
6. **Filters**: YAML files that define how to extract and transform data from raw events.

The EventProcessor uses a plugin architecture where plugins are separate processes that communicate with the
EventProcessor via gRPC over Unix sockets. This allows for a flexible and extensible system where new functionality can
be added without modifying the core engine.

## Components

### EventProcessor

The EventProcessor is the core engine that processes security events. It loads plugins, routes events to the appropriate
plugins, and manages the overall flow of data through the system.

### EventProcessor Plugins

The EventProcessor has several built-in plugins for different tasks:

- **Input Plugins**: Collect or receive logs from external sources (http-input, grpc-input).
- **Parsing Plugins**: Extract and enrich data and transform logs (add, cast, csv, delete, expand, grok, json, kv,
  reformat, rename, trim).
- **Analysis Plugins**: Process logs to detect security attacks (analysis).

### UTMStack Plugins

UTMStack has several plugins that integrate with the EventProcessor:

- **Input Plugins**: Collect events from various sources (aws, azure, bitdefender, gcp, o365, sophos, inputs).
- **Parsing Plugins**: Extract and enrich data and transform logs (geolocation).
- **Analysis Plugins**: Process and analyze events (events).
- **Correlation Plugins**: Detect relationships between alerts (alerts, soc-ai).
- **Notification Plugins**: Send notifications or statistics to internal and external systems (stats).
- **Sidecar Plugins**: Background task plugins with multiple purposes, like maintaining the system healthy or managing
  configurations (config).

### go-sdk

The go-sdk provides common functionality for both the EventProcessor and plugins. It defines the interfaces and types
used for communication between components.

Key files in the go-sdk:

- `plugins/plugins.proto`: Defines the Protocol Buffers messages and services used by the plugin system.
- `plugins/config.go`: Provides configuration functionality for plugins.
- `plugins/cel.go`: Provides Common Expression Language (CEL) functionality for rules.

## Implementing Rules

Rules are YAML files that define how to analyze events to detect security threats. They're used by the analysis plugin
to generate alerts when specific conditions are met.

### Rule Structure

A rule is defined as a YAML object with the following fields:

```yaml
- id: 1                           # Unique identifier for the rule
  dataTypes: # Types of data this rule applies to
    - google
  name: Hello                     # Name of the rule
  impact: # Impact information
    confidentiality: 0            # Impact on confidentiality (0-5)
    integrity: 0                  # Impact on integrity (0-5)
    availability: 3               # Impact on availability (0-5)
  category: Testing Category      # Category of the rule
  technique: Testing Technique    # Technique used by the threat
  adversary: origin               # Which side is considered the adversary (origin or target)
  references: # External references
    - https://quantfall.com
  description: This is a testing rule.  # Description of the rule
  where: # Conditions for when the rule applies
    variables: # Variables to extract from the event
      - get: origin.geolocation.country  # Path to the value in the event
        as: country                      # Name of the variable
        ofType: "string"                 # Type of the variable (required)
    expression: country_ok && country == "United States"  # Expression to evaluate
  afterEvents: # Additional events to search for
    - indexPattern: v11-log-*     # Index pattern to search in
      with: # Conditions for the search
        - field: origin.ip.keyword  # Field to match
          operator: eq              # Operator (eq, neq)
          value: '{{origin.ip}}'    # Value to match (can use variables from the event)
      within: now-12h             # Time window for the search
      count: 1                    # Number of events to match
  deduplicateBy: # Fields used for deduplication
    - adversary.ip
    - adversary.country
```

### Rule Fields

- **id**: A unique identifier for the rule.
- **dataTypes**: An array of data types that this rule applies to. The rule will only be evaluated for events with these
  data types.
- **name**: The name of the rule.
- **impact**: The impact of the threat detected by this rule, with scores for confidentiality, integrity, and
  availability.
- **category**: The category of the rule.
- **technique**: The technique used by the threat.
- **adversary**: Which side is considered the adversary (origin or target).
- **references**: An array of external references for more information about the threat.
- **description**: A description of the rule.
- **where**: Conditions for when the rule applies.
    - **variables**: Variables to extract from the event.
        - **get**: The path to the value in the event.
        - **as**: The name of the variable.
        - **ofType**: The type of the variable (required). Possible values include "string", "int", "double", "bool", "
          bytes", "uint", "timestamp", "duration", "type", "null", "any", and various list and map types.
    - **expression**: An expression to evaluate using the variables.
- **afterEvents**: Additional events to search for.
    - **indexPattern**: The index pattern to search in.
    - **with**: Conditions for the search.
        - **field**: The field to match.
        - **operator**: The operator to use for matching. Possible values:
            - **eq**: Equality operator. Matches events where the field equals the value.
            - **neq**: Not equal operator. Matches events where the field does not equal the value.
        - **value**: The value to match (can use variables from the event using the `{{field.path}}` syntax).
    - **within**: The time window for the search.
    - **count**: The number of events to match.
- **deduplicateBy**: Fields used for deduplication of alerts.

### Rule Evaluation

When an event is received, the analysis plugin evaluates all rules that apply to the event's data type. For each rule:

1. The variables are extracted from the event.
2. The expression is evaluated using the variables.
3. If the expression evaluates to true, the afterEvents searches are performed.
4. If all conditions are met, an alert is generated.

## Implementing Filters

Filters are YAML files that define how to extract and transform data from raw events. They are used by the parsing
plugin to convert raw events into a standardized format that can be used by the EventProcessor and analyzed by rules.

### Filter Structure

A filter is defined as a YAML object with the following structure:

```yaml
pipeline:
  - dataTypes: # Types of data this filter applies to
      - wineventlog
    steps: # Processing steps to apply
      - json: # Parse the raw data as JSON
          source: raw
      - rename: # Rename fields
          from:
            - log.host.ip
          to: origin.ip
      # More steps...
```

### Filter Steps

Filters can include various types of steps for processing events. Each step type serves a specific purpose in the data
transformation pipeline:

#### 1. **json** - Parse JSON Data

Parses a field containing JSON data into structured log fields.

**Required Fields:**

- `source`: The field containing JSON data

**Example:**

```yaml
- json:
    source: raw
```

#### 2. **rename** - Rename Fields

Renames existing fields to match standardized field naming conventions.

**Required Fields:**

- `from`: Array of source field names
- `to`: Target field name

**Example:**

```yaml
- rename:
    from:
      - log.host.ip
      - log.source.ip
    to: origin.ip
```

#### 3. **cast** - Type Conversion

Converts field values to specified data types.

**Required Fields:**

- `fields`: Array of field names to convert
- `to`: Target data type (`int`, `float`, `float64`, `string`, `[]string`, etc.)

**Example:**

```yaml
- cast:
    fields:
      - origin.port
      - statusCode
    to: int

- cast:
    fields:
      - log.local.ips
    to: '[]string'
```

#### 4. **delete** - Remove Fields

Removes specified fields from the log structure.

**Required Fields:**

- `fields`: Array of field names to remove

**Optional Fields:**

- `where`: Conditional logic for when to apply the deletion

**Example:**

```yaml
- delete:
    fields:
      - log.method
      - log.service
      - log.metadata
    where: has(action)
```

#### 5. **grok** - Pattern-based Parsing

Extracts structured data from unstructured text using pattern matching. Each grok step uses a list of patterns that are
applied sequentially to extract multiple fields from the source.

**Required Fields:**

- `patterns`: Array of pattern definitions, each with `fieldName` and `pattern`
- `source`: Source field to parse

**Note:** Use `fieldName` (camelCase) in pattern definitions. Some legacy filters may use `field_name` but `fieldName`
is the standard.

**Optional Fields:**

- `where`: Conditional logic for when to apply parsing

**Pattern Structure:**

- Each pattern in the list defines one field to extract
- Patterns are applied in sequence to parse complex log formats
- Multiple grok steps can be used for different parsing stages

**Built-in Patterns:**

- `{{.ipv4}}`, `{{.ipv6}}` - IP addresses
- `{{.integer}}`, `{{.word}}` - Numbers and words
- `{{.data}}`, `{{.greedy}}` - Generic data patterns
- `{{.time}}`, `{{.year}}`, `{{.monthNumber}}`, `{{.monthDay}}` - Time patterns
- `{{.hostname}}`, `{{.day}}`, `{{.monthName}}` - Host and date patterns

**Examples:**

Basic Apache log parsing:

```yaml
- grok:
    patterns:
      - fieldName: origin.ip
        pattern: '{{.ipv4}}|{{.ipv6}}'
      - fieldName: origin.user
        pattern: '{{.word}}|(-)'
      - fieldName: deviceTime
        pattern: '\[{{.data}}\]'
      - fieldName: log.request
        pattern: '\"{{.data}}\"'
      - fieldName: log.statusCode
        pattern: '{{.integer}}'
    source: log.message
```

Complex Cisco ASA parsing:

```yaml
- grok:
    patterns:
      - fieldName: log.ciscoTime
        pattern: '({{.day}}\s)?{{.monthName}}\s{{.monthDay}}\s{{.year}}\s{{.time}}'
      - fieldName: log.localIp
        pattern: '{{.ipv4}}|{{.ipv6}}|{{.hostname}}'
      - fieldName: log.asaHeader
        pattern: '{{.data}}ASA-'
      - fieldName: log.severity
        pattern: '{{.integer}}'
      - fieldName: log.messageId
        pattern: '-{{.integer}}'
      - fieldName: log.msg
        pattern: '{{.greedy}}'
    source: raw
```

Parsing with port extraction:

```yaml
- grok:
    patterns:
      - fieldName: origin.ip
        pattern: '({{.ipv4}}|{{.ipv6}})'
      - fieldName: origin.port
        pattern: '/{{.integer}}'
      - fieldName: target.ip
        pattern: '({{.ipv4}}|{{.ipv6}})'
      - fieldName: target.port
        pattern: '/{{.integer}}'
    source: log.connectionInfo
```

#### 6. **kv** - Key-Value Parsing

Parses key-value formatted data into structured fields.

**Required Fields:**

- `fieldSplit`: Character(s) that separate key-value pairs
- `valueSplit`: Character(s) that separate keys from values
- `source`: Source field containing key-value data

**Optional Fields:**

- `where`: Conditional logic for when to apply parsing

**Example:**

```yaml
- kv:
    fieldSplit: " "
    valueSplit: "="
    source: raw
```

#### 7. **trim** - String Trimming

Removes specified characters from the beginning or end of field values.

**Required Fields:**

- `function`: Trim operation (`prefix`, `suffix`)
- `substring`: Character(s) to remove
- `fields`: Array of fields to trim

**Optional Fields:**

- `where`: Conditional logic for when to apply trimming

**Example:**

```yaml
- trim:
    function: prefix
    substring: '['
    fields:
      - deviceTime
      - log.severityLabel

- trim:
    function: suffix
    substring: ':'
    fields:
      - origin.ip
```

#### 8. **add** - Add New Fields

Adds new fields with specified values, often used for field enrichment and normalization.

**Required Fields:**

- `function`: Add function type (`string`)
- `params`: Parameters including `key` and `value`

**Optional Fields:**

- `where`: Conditional logic for when to add the field

**Example:**

```yaml
- add:
    function: 'string'
    params:
      key: actionResult
      value: 'accepted'
    where: has(statusCode) && (statusCode >= 200 && statusCode <= 299)

- add:
    function: 'string'
    params:
      key: action
      value: 'get'
    where:
      variables: has(log.method) && log.method == "GET"
```

#### 9. **reformat** - Field Reformatting

Reformats field values, particularly useful for timestamp conversion.

**Required Fields:**

- `fields`: Array of fields to reformat
- `function`: Reformatting function (`time`)
- `fromFormat`: Source format pattern
- `toFormat`: Target format pattern

**Example:**

```yaml
- reformat:
    fields:
      - deviceTime
    function: time
    fromFormat: '14/Feb/2022:15:40:53 -0500'
    toFormat: '2024-09-23T15:57:40.338364445Z'
```

#### 10. **expand** - Field Expansion

Expands complex nested fields or structured data into individual fields.

**Required Fields:**

- `source`: Source field containing data to expand
- `to`: Target field name for expanded data

**Optional Fields:**

- `where`: Conditional logic for when to apply expansion

**Example:**

```yaml
- expand:
    source: log.jsonPayload.structuredRdata
    to: log.jsonPayloadStructuredRdata
    where: has(log.jsonPayload.structuredRdata)
```

#### 11. **csv** - CSV Data Parsing

Parses CSV-formatted data into structured fields using defined column mappings.

**Required Fields:**

- `source`: Source field containing CSV data
- `separator`: Character used to separate CSV values (typically `","`)
- `headers`: Array of field names for CSV columns

**Optional Fields:**

- `where`: Conditional logic for when to apply CSV parsing

**Note:** The protocol buffer definition only supports `headers`, not `columns`. Some filters may incorrectly use
`columns` but `headers` is the standard.

**Example:**

```yaml
- csv:
    source: log.csvMsg
    separator: ","
    headers:
      - log.ruleNumber
      - log.subRuleNumber
      - log.anchor
      - log.tracker
      - log.realInterface
      - log.reason
      - log.action
```

#### 12. **dynamic** - Dynamic Plugin Integration

Calls external plugins for specialized processing like geolocation enrichment.

**Required Fields:**

- `plugin`: Plugin identifier
- `params`: Plugin-specific parameters

**Optional Fields:**

- `where`: Conditional logic for when to apply the plugin

**Example:**

```yaml
- dynamic:
    plugin: com.utmstack.geolocation
    params:
      source: origin.ip
      destination: origin.geolocation
    where: has(origin.ip)
```

### Conditional Logic in Steps

All steps support conditional logic using the `where` clause:

**Structure:**

```yaml
where: "conditional_expression"
```

**Supported Data Types:**

- `string`, `int`, `float`, `bool`
- `[]string`, `[]int` (arrays)
- Custom types as needed

**Expression Examples:**

- `has(origin.ip)` - Check if variable exists and is valid
- `statusCode >= 200 && statusCode <= 299` - Range checking
- `log.method == "GET"` - String comparison
- `severity.contains("error")` - String operations

### Complete Filter Example

Here's a comprehensive example showing multiple step types working together (based on Apache access log processing):

```yaml
pipeline:
  - dataTypes:
      - apache
    steps:
      # 1. Parse JSON structure from raw input
      - json:
          source: raw

      # 2. Rename fields to standardized mapping
      - rename:
          from:
            - log.host.hostname
          to: origin.host
      - rename:
          from:
            - log.host.ip
          to: log.local.ips

      # 3. Type conversion for array fields
      - cast:
          to: '[]string'
          fields:
            - log.local.ips

      # 4. Parse Apache Common Log Format using grok
      - grok:
          patterns:
            - fieldName: origin.ip
              pattern: '{{.ipv4}}|{{.ipv6}}'
            - fieldName: origin.user
              pattern: '{{.word}}|(-)'
            - fieldName: deviceTime
              pattern: '\[{{.data}}\]'
            - fieldName: log.request
              pattern: '\"{{.data}}\"'
            - fieldName: log.statusCode
              pattern: '{{.integer}}'
            - fieldName: origin.bytesReceived
              pattern: '{{.integer}}|(-)'
          source: log.message

      # 5. Clean up parsed fields
      - trim:
          function: prefix
          substring: '['
          fields:
            - deviceTime
      - trim:
          function: suffix
          substring: ']'
          fields:
            - deviceTime

      # 6. Parse HTTP request components
      - grok:
          patterns:
            - fieldName: log.method
              pattern: '{{.word}}'
            - fieldName: origin.path
              pattern: '(.*)\s+'
            - fieldName: protocol
              pattern: '{{.greedy}}'
          source: log.request

      # 7. Add geolocation data using dynamic plugin
      - dynamic:
          plugin: com.utmstack.geolocation
          params:
            source: origin.ip
            destination: origin.geolocation
          where: has(origin.ip)

      # 8. Normalize HTTP methods to standardized actions
      - add:
          function: 'string'
          params:
            key: action
            value: 'get'
          where: has(log.method) && log.method == "GET"

      # 9. Add result classification based on status codes
      - add:
          function: 'string'
          params:
            key: actionResult
            value: 'accepted'
          where: has(log.statusCode) && (log.statusCode >= 200 && log.statusCode <= 299)

      # 10. Convert numeric fields
      - cast:
          fields:
            - log.statusCode
            - origin.bytesReceived
          to: int

      # 11. Reformat timestamp
      - reformat:
          fields:
            - deviceTime
          function: time
          fromFormat: '14/Feb/2022:15:40:53 -0500'
          toFormat: '2024-09-23T15:57:40.338364445Z'

      # 12. Clean up temporary fields
      - delete:
          fields:
            - log.method
            - log.service
            - log.agent
          where: has(action)
```

This example demonstrates:

- **Sequential processing**: Steps are applied in order
- **Field transformation**: Renaming, casting, and reformatting
- **Pattern extraction**: Using grok for complex log parsing
- **Conditional logic**: Adding fields based on conditions
- **External integration**: Geolocation enrichment via dynamic plugin
- **Data cleanup**: Trimming and deleting unnecessary fields

### Filter Evaluation

When an event is received, the parsing plugin selects the appropriate filter based on the event's data type. The filter
is then applied to the event, transforming it according to the defined steps. Each step modifies the event structure,
preparing it for analysis and correlation by downstream plugins.

## Best Practices

### Rule Development

1. **Start Simple**: Begin with simple rules that match specific patterns, then refine them as needed.
2. **Test Thoroughly**: Test rules with a variety of events to ensure they work as expected.
3. **Use Variables**: Use variables to make rules more readable and maintainable.
4. **Document Rules**: Include a clear description and references in each rule.
5. **Consider Performance**: Complex rules can impact performance, so optimize them as needed.

### Filter Development

1. **Standardize Field Names**: Use consistent field names across all filters.
2. **Remove Unnecessary Fields**: Delete fields that are not needed for analysis to reduce storage requirements.
3. **Handle Edge Cases**: Consider how to handle missing or malformed data.
4. **Document Filters**: Include comments to explain the purpose of each step.
5. **Test with Real Data**: Test filters with real data to ensure they work as expected.

## Development Workflow

This section provides a step-by-step guide for developing and implementing new rules and filters.

### Rule Development Workflow

1. **Identify the Security Threat**: Determine what security threat you want to detect.
2. **Understand the Data**: Examine the events that would indicate this threat.
3. **Create a Rule File**: Create a new YAML file in the rules directory.
4. **Define Basic Metadata**: Set the id, name, description, and other metadata.
5. **Define Data Types**: Specify which data types this rule applies to.
6. **Define Impact**: Set the confidentiality, integrity, and availability impact scores.
7. **Define Where Conditions**: Create variables and an expression to identify events of interest.
8. **Define After Events**: If needed, specify additional events to search for.
9. **Define Deduplication**: Specify fields to use for deduplicating alerts.
10. **Test the Rule**: Deploy the rule and test it with sample events.
11. **Refine the Rule**: Adjust the rule based on testing results.
12. **Document the Rule**: Add comments and references to explain the rule.

### Filter Development Workflow

1. **Identify the Data Source**: Determine what data source you want to process.
2. **Understand the Raw Format**: Examine the raw events from this source.
3. **Create a Filter File**: Create a new YAML file in the appropriate filters directory.
4. **Define Data Types**: Specify which data types this filter applies to.
5. **Define Parsing Steps**: Add steps to parse and transform the raw data.
6. **Test the Filter**: Deploy the filter and test it with sample events.
7. **Refine the Filter**: Adjust the filter based on testing results.
8. **Document the Filter**: Add comments to explain the filter.

## Advanced Features

This section covers advanced features of the EventProcessor and UTMStack ecosystem that can be used to create more
sophisticated rules and filters.

### Advanced Rule Features

#### Complex Expressions

The `where.expression` field in rules supports complex expressions using the Common Expression Language (CEL). CEL is a
powerful expression language that allows for complex logic, including:

- **Logical Operators**: `&&` (AND), `||` (OR), `!` (NOT)
- **Comparison Operators**: `==`, `!=`, `<`, `<=`, `>`, `>=`
- **String Operations**: `startsWith()`, `endsWith()`, `contains()`
- **Array Operations**: `in`, `size()`
- **Mathematical Operations**: `+`, `-`, `*`, `/`, `%`

Example of a complex expression:

```yaml
expression: has(origin.country) && !(origin.country in ["United States", "Canada", "United Kingdom"]) && (origin.user != "" && origin.user.startsWith("admin"))
```

#### Nested AfterEvents

The `afterEvents` field in rules supports nested searches using the `or` field. This allows for more complex correlation
logic:

```yaml
afterEvents:
  - indexPattern: v11-log-*
    with:
      - field: origin.ip.keyword
        operator: eq
        value: '{{origin.ip}}'
    within: now-12h
    count: 1
    or:
      - indexPattern: v11-alert-*
        with:
          - field: adversary.ip.keyword
            operator: eq
            value: '{{origin.ip}}'
        within: now-24h
        count: 2
```

In this example, the rule will match if either:

1. There is at least 1 event in the `v11-log-*` index with the same origin IP within the last 12 hours, OR
2. There are at least 2 alerts in the `v11-alert-*` index with the same adversary IP within the last 24 hours.

#### Dynamic Values

Rule fields can use dynamic values from the event using the `{{field.path}}` syntax. This is particularly useful in the
`afterEvents` section:

```yaml
afterEvents:
  - indexPattern: v11-log-*
    with:
      - field: origin.user.keyword
        operator: eq
        value: '{{origin.user}}'
      - field: origin.ip.keyword
        operator: neq
        value: '{{origin.ip}}'
    within: now-24h
    count: 3
```

This example searches for events with the same user but a different IP address, which could indicate a compromised
account.

### Advanced Filter Features

#### Multi-Stage Pipelines

Filters can include multiple stages in the pipeline, each with its own set of steps:

```yaml
pipeline:
  - dataTypes:
      - wineventlog
    steps:
      - json:
          source: raw
      - rename:
          from:
            - log.host.ip
          to: origin.ip
  - dataTypes:
      - wineventlog
    steps:
      - grok:
          source: message
          pattern: "%{WORD:action} %{IP:target.ip}"
          target: parsed
```

This allows for more modular and maintainable filters, especially for complex data sources.

## Integration with Other Systems

The EventProcessor and UTMStack ecosystem can integrate with various other systems to enhance its capabilities.

### Integration with Threat Intelligence Platforms

UTMStack can integrate with threat intelligence platforms to enrich events with threat intelligence data:

1. **ThreatWinds**: UTMStack has native integration with ThreatWinds for threat intelligence.
2. **MISP**: UTMStack can integrate with MISP to consume threat intelligence feeds via third party plugins.
3. **AlienVault OTX**: UTMStack can integrate with AlienVault OTX to consume threat intelligence feeds via third party
   plugins.

### Integration with Ticketing Systems

UTMStack can integrate with ticketing systems to create tickets for alerts:

1. **JIRA**: UTMStack can create JIRA tickets for alerts using the JIRA API via third party plugins.
2. **ServiceNow**: UTMStack can create ServiceNow incidents for alerts using the ServiceNow API via third party plugins.
3. **GitHub Issues**: UTMStack can create GitHub issues for alerts using the GitHub API via third party plugins.

### Integration with Communication Platforms

UTMStack can integrate with communication platforms to send notifications for alerts:

1. **Email**: UTMStack can send email notifications for alerts.
2. **Slack**: UTMStack can send Slack notifications for alerts using the Slack API via third party plugins.
3. **Microsoft Teams**: UTMStack can send Microsoft Teams notifications for alerts using the Microsoft Teams API via
   third party plugins.

## Performance Optimization

This section provides guidance on optimizing the performance of the EventProcessor and UTMStack ecosystem.

### Rule Optimization

1. **Limit Data Types**: Specify only the data types that the rule applies to. This reduces the number of events that
   need to be evaluated.
2. **Use Efficient Expressions**: Use efficient expressions in the `where` field. Avoid complex expressions
   that require a lot of processing.
3. **Limit AfterEvents Searches**: Limit the number of `afterEvents` searches and the time window for each search. This
   reduces the load on the search engine.
4. **Use Deduplication**: Use the `deduplicateBy` field to prevent alert fatigue.

### Filter Optimization

1. **Limit Data Types**: Specify only the data types that the filter applies to. This reduces the number of events that
   need to be processed.
2. **Use Efficient Steps**: Use efficient steps in the filter pipeline. Avoid complex steps that require a lot of
   processing.
3. **Remove Unnecessary Fields**: Remove fields that are not needed for analysis to reduce storage requirements.
4. **Use Conditional Steps**: Use conditional steps to apply different processing based on the event type. This can
   reduce the number of steps that need to be applied to each event.

### System Optimization

1. **Hardware Resources**: Ensure that the system has sufficient hardware resources (CPU, memory, disk) to handle the
   expected event volume.
2. **Cluster Configuration**: Configure the OpenSearch cluster with the appropriate number of nodes, shards, and
   replicas for the expected event volume.
3. **JVM Settings**: Configure the JVM settings for OpenSearch to optimize memory usage.
4. **Network Configuration**: Ensure that the network configuration is optimized for the expected event volume.

## Troubleshooting

### Common Issues

1. **Rule Not Triggering**: Check that the event matches the dataTypes and where conditions.
2. **Filter Not Processing**: Check that the event matches the dataTypes and that the filter steps are correct.
3. **Missing Fields**: Check that the fields referenced in rules and filters exist in the events.
4. **Performance Issues**: Check for complex rules or filters that may be impacting performance.

### Debugging

1. **Check Logs**: Look for error messages in the EventProcessor and plugin logs.
2. **Test Rules Individually**: Test rules one at a time to isolate issues.
3. **Validate YAML**: Ensure that rule and filter YAML files are valid.
4. **Check Field Names**: Verify that field names in rules and filters match the actual field names in events.
5. **Use Test Events**: Create test events that should trigger your rules and verify they work as expected.

## Real-World Use Cases

This section provides real-world examples of how the EventProcessor and UTMStack ecosystem can be used to solve security
challenges.

### Detecting Brute Force Attacks

A common security challenge is detecting brute force attacks against authentication systems. Here's how you can use the
EventProcessor and UTMStack to detect such attacks:

1. **Create a Filter**: Create a filter that extracts relevant information from authentication logs, such as the source
   IP, username, and authentication result.

```yaml
pipeline:
  - dataTypes:
      - auth_logs
    steps:
      - json:
          source: raw
      - rename:
          from:
            - log.source.ip
          to: origin.ip
      - rename:
          from:
            - log.auth.username
          to: origin.user
      - rename:
          from:
            - log.auth.result
          to: actionResult
```

2. **Create a Rule**: Create a rule that detects multiple failed authentication attempts from the same IP address.

```yaml
- id: 201
  dataTypes:
    - auth_logs
  name: Brute Force Attack Detection
  impact:
    confidentiality: 4
    integrity: 3
    availability: 2
  category: Authentication
  technique: Brute Force
  adversary: origin
  references:
    - https://attack.mitre.org/techniques/T1110/
  description: Detects multiple failed authentication attempts from the same IP address.
  where: has(origin.ip) && actionResult == "failure"
  afterEvents:
    - indexPattern: v11-log-auth_logs
      with:
        - field: origin.ip.keyword
          operator: eq
          value: '{{origin.ip}}'
        - field: actionResult.keyword
          operator: eq
          value: 'failure'
      within: now-1h
      count: 5
  deduplicateBy:
    - origin.ip
```

This rule will generate an alert if there are at least 5 failed authentication attempts from the same IP address within
the last hour.

### Detecting Data Exfiltration

Another common security challenge is detecting data exfiltration. Here's how you can use the EventProcessor and UTMStack
to detect unusual data transfers:

1. **Create a Filter**: Create a filter that extracts relevant information from network logs, such as the source IP,
   destination IP, and data transfer size.

```yaml
pipeline:
  - dataTypes:
      - network_logs
    steps:
      - json:
          source: raw
      - rename:
          from:
            - log.source.ip
          to: origin.ip
      - rename:
          from:
            - log.destination.ip
          to: target.ip
      - rename:
          from:
            - log.network.bytes
          to: origin.bytesSent
```

2. **Create a Rule**: Create a rule that detects large data transfers to unusual destinations.

```yaml
- id: 202
  dataTypes:
    - network_logs
  name: Data Exfiltration Detection
  impact:
    confidentiality: 5
    integrity: 2
    availability: 1
  category: Exfiltration
  technique: Data Transfer
  adversary: origin
  references:
    - https://attack.mitre.org/techniques/T1048/
  description: Detects large data transfers to unusual destinations.
  where: has(origin.ip) && has(target.ip) && has(origin.bytesSent) origin.bytesSent && > 10000000 && !(target.ip in ["10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"])
  afterEvents:
    - indexPattern: v11-log-network_logs
      with:
        - field: origin.ip.keyword
          operator: eq
          value: '{{origin.ip}}'
        - field: target.ip.keyword
          operator: eq
          value: '{{target.ip}}'
      within: now-24h
      count: 1
  deduplicateBy:
    - origin.ip
    - target.ip
```

This rule will generate an alert if there is a data transfer larger than 10MB to a destination outside the internal
network.

### Detecting Insider Threats

Insider threats can be difficult to detect. Here's how you can use the EventProcessor and UTMStack to detect unusual
user behavior that might indicate an insider threat:

1. **Create a Filter**: Create a filter that extracts relevant information from user activity logs, such as the
   username, action, and time.

```yaml
pipeline:
  - dataTypes:
      - user_activity
    steps:
      - json:
          source: raw
      - rename:
          from:
            - log.user.name
          to: origin.user
      - rename:
          from:
            - log.activity.action
          to: action
      - rename:
          from:
            - log.activity.time
          to: deviceTime
```

2. **Create a Rule**: Create a rule that detects unusual user activity outside normal working hours.

```yaml
- id: 203
  dataTypes:
    - user_activity
  name: Unusual User Activity
  impact:
    confidentiality: 3
    integrity: 4
    availability: 2
  category: Insider Threat
  technique: Unusual Activity
  adversary: origin
  references:
    - https://attack.mitre.org/techniques/T1078/
  description: Detects unusual user activity outside normal working hours.
  where: has(origin.user) && has(deviceTime) && has(action) && (time.hour < 8 || time.hour > 18) && action in ["file_access", "database_query", "admin_action"]
  afterEvents:
    - indexPattern: v11-log-user_activity
      with:
        - field: origin.user.keyword
          operator: eq
          value: '{{origin.user}}'
        - field: action.keyword
          operator: eq
          value: '{{action}}'
      within: now-7d
      count: 1
  deduplicateBy:
    - origin.user
    - action
```

This rule will generate an alert if a user performs sensitive actions outside normal working hours (8 AM to 6 PM).

## Custom Plugin Development

The EventProcessor and UTMStack ecosystem is designed to be extensible through plugins. This section provides guidance
on developing custom plugins to extend the functionality of the system.

### Plugin Types

There are five main types of plugins that you can develop for the EventProcessor, each serving specific purposes in the
security event processing pipeline:

1. **Input Plugins**: Collect or receive logs from external sources (e.g., syslog, APIs, files).
2. **Parsing Plugins**: Extract and transform logs to match a defined mapping, enabling analysis/correlation and
   allowing UTMStack users to explore logs in the Logs Explorer and visualize aggregated data in dashboards.
3. **Analysis Plugins**: Process logs to detect security attacks and generate alerts.
4. **Correlation Plugins**: Detect when different alerts are associated or have correlation relationships.
5. **Notification Plugins**: Send notifications to UTMStack and external systems (e.g., Slack channels, ticketing
   systems, email).

### Plugin Structure

A plugin is a separate process that communicates with the EventProcessor via gRPC over Unix sockets. Each plugin type
implements a specific gRPC service interface:

1. **Input Plugins**: Implement the `Engine` service (bidirectional streaming) to send logs to the EventProcessor
2. **Parsing Plugins**: Implement the `Parsing` service (unary RPC) with `ParseLog(Transform) -> Draft` method
3. **Analysis Plugins**: Implement the `Analysis` service (server streaming) with `Analyze(Event) -> stream Alert`
   method
4. **Correlation Plugins**: Implement the `Correlation` service to process and correlate alerts
5. **Notification Plugins**: Implement the `Notification` service to send notifications to external systems

The basic structure of a plugin includes:

1. **Main Package**: The entry point for the plugin.
2. **gRPC Server**: Implements the appropriate gRPC service interface based on plugin type.
3. **Plugin Logic**: Implements the specific functionality for the plugin's purpose in the pipeline.
4. **Unix Socket**: Creates and listens on a named Unix socket for communication with the EventProcessor.

### Developing an Input Plugin

Input plugins collect or receive logs from external sources and send them to the EventProcessor. They implement the
`Engine` service interface.

Here's a step-by-step guide for developing a custom input plugin:

1. **Create a New Go Module**: Create a new Go module for your plugin.

```bash
mkdir my-input-plugin
cd my-input-plugin
go mod init github.com/myorg/my-input-plugin
```

2. **Add Dependencies**: Add the go-sdk as a dependency.

```bash
go get github.com/threatwinds/go-sdk
```

3. **Create the Main Package**: Create a main.go file with the basic structure of the input plugin.

```go
package main

import (
	"context"
	"os"
	"time"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	"github.com/threatwinds/go-sdk/utils"
	"google.golang.org/grpc"
)

func main() {
	// Try to set up the connection
	socketsFolder, err := utils.MkdirJoin(plugins.WorkDir, "sockets")
	if err != nil {
		_ = catcher.Error("cannot create sockets folder", err, nil)
		time.Sleep(5 * time.Second) // Wait before retrying
		os.Exit(1)
	}

	socket := socketsFolder.FileJoin("engine_server.sock")

	conn, err := grpc.NewClient(
		fmt.Sprintf("unix://%s", socket),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		_ = catcher.Error("cannot connect to engine server", err, map[string]any{})
		time.Sleep(5 * time.Second) // Wait before retrying
		os.Exit(1)
	}

	defer conn.Close()

	client := plugins.NewEngineClient(conn)

	inputStreamClient, err := client.Input(context.Background())
	if err != nil {
		_ = catcher.Error("cannot create input client", err, map[string]any{})
		time.Sleep(5 * time.Second) // Wait before retrying
		os.Exit(1)
	}

	var logsChannel = make(chan string, 100)

	// Start collecting data from an external source
	go collectData(logsChannel)

	// Start sending data to the EventProcessor
	go sendData(inputStreamClient, logsChannel)

	// Keep the plugin running
	select {}
}

func sendData(inputStreamClient grpc.BidiStreamingClient[plugins.Log, plugins.Ack], logs chan string) {
	for {
		log := <-logs
		// Implement bidirectional streaming logic here
	}
}

func collectData(logs chan string) {
	// Implement data collection logic here
	logs <- "collected data"
}
```

4. **Build the Plugin**: Build the plugin as a standalone executable.

```bash
go build -o my-input-plugin
```

5. **Deploy the Plugin**: Deploy the plugin to the plugins directory of the EventProcessor.

```bash
cp my-input-plugin /path/to/eventprocessor/plugins/
```

### Developing a Parsing Plugin

Parsing plugins extract and transform logs to match a defined mapping that enables analysis/correlation and allows
UTMStack users to explore logs in the Logs Explorer. They implement the `Parsing` service interface.

Here's a step-by-step guide for developing a custom parsing plugin:

1. **Create a New Go Module**: Create a new Go module for your plugin.

```bash
mkdir my-parsing-plugin
cd my-parsing-plugin
go mod init github.com/myorg/my-parsing-plugin
```

2. **Add Dependencies**: Add the go-sdk as a dependency.

```bash
go get github.com/threatwinds/go-sdk
```

3. **Create the Main Package**: Create a main.go file with the basic structure of the parsing plugin.

```go
package main

import (
	"context"
	"net"
	"os"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	"github.com/threatwinds/go-sdk/utils"
	"google.golang.org/grpc"
)

type parsingServer struct {
	plugins.UnimplementedParsingServer
}

func main() {
	// Create socket directory
	filePath, err := utils.MkdirJoin(plugins.WorkDir, "sockets")
	if err != nil {
		_ = catcher.Error("cannot create socket directory", err, nil)
		os.Exit(1)
	}

	// Create socket path for parsing plugin
	socketPath := filePath.FileJoin("my_parsing.sock")
	_ = os.Remove(socketPath)

	// Resolve Unix address
	unixAddress, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		_ = catcher.Error("cannot resolve unix address", err, nil)
		os.Exit(1)
	}

	// Listen on Unix socket
	listener, err := net.ListenUnix("unix", unixAddress)
	if err != nil {
		_ = catcher.Error("cannot listen to unix socket", err, nil)
		os.Exit(1)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	plugins.RegisterParsingServer(grpcServer, &parsingServer{})

	// Serve gRPC
	if err := grpcServer.Serve(listener); err != nil {
		_ = catcher.Error("cannot serve grpc", err, nil)
		os.Exit(1)
	}
}

func (p *parsingServer) ParseLog(ctx context.Context, transform *plugins.Transform) (*plugins.Draft, error) {
	// Implement your parsing logic here
	// Note: transform.Draft contains log as string, not structured data
	// The actual log parsing is handled by the EventProcessor

	return transform.Draft, nil
}
```

4. **Implement the Parsing Logic**: The `ParseLog` method receives a `Transform` containing the current draft and step
   configuration, and returns a modified `Draft`.

5. **Build the Plugin**: Build the plugin as a standalone executable.

```bash
go build -o my-parsing-plugin
```

6. **Deploy the Plugin**: Deploy the plugin to the plugins directory of the EventProcessor.

```bash
cp my-parsing-plugin /path/to/eventprocessor/plugins/
```

### Developing an Analysis Plugin

Analysis plugins process logs to detect security attacks and generate alerts. They implement the `Analysis` service
interface.

```go
package main

import (
	"net"
	"os"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	"github.com/threatwinds/go-sdk/utils"
	"google.golang.org/grpc"
)

type analysisServer struct {
	plugins.UnimplementedAnalysisServer
}

func main() {
	// Create socket directory
	filePath, err := utils.MkdirJoin(plugins.WorkDir, "sockets")
	if err != nil {
		_ = catcher.Error("cannot create socket directory", err, nil)
		os.Exit(1)
	}

	// Create socket path for analysis plugin
	socketPath := filePath.FileJoin("my_analysis.sock")
	_ = os.Remove(socketPath)

	// Resolve Unix address
	unixAddress, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		_ = catcher.Error("cannot resolve unix address", err, nil)
		os.Exit(1)
	}

	// Listen on Unix socket
	listener, err := net.ListenUnix("unix", unixAddress)
	if err != nil {
		_ = catcher.Error("cannot listen to unix socket", err, nil)
		os.Exit(1)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	plugins.RegisterAnalysisServer(grpcServer, &analysisServer{})

	// Serve gRPC
	if err := grpcServer.Serve(listener); err != nil {
		_ = catcher.Error("cannot serve grpc", err, nil)
		os.Exit(1)
	}
}

func (a *analysisServer) Analyze(event *plugins.Event, srv grpc.ServerStreamingServer[plugins.Alert]) error {
	// Implement your analysis logic here
	// Example: Detect failed login attempts

	return nil // No alert generated
}
```

### Developing a Correlation Plugin

Correlation plugins detect when different alerts are associated or have correlation relationships.

```go
package main

import (
	"context"
	"net"
	"os"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	"github.com/threatwinds/go-sdk/utils"
	"google.golang.org/grpc"
)

type correlationServer struct {
	plugins.UnimplementedCorrelationServer
}

func main() {
	// Create socket directory
	filePath, err := utils.MkdirJoin(plugins.WorkDir, "sockets")
	if err != nil {
		_ = catcher.Error("cannot create socket directory", err, nil)
		os.Exit(1)
	}

	// Create socket path for correlation plugin
	socketPath := filePath.FileJoin("my_correlation.sock")
	_ = os.Remove(socketPath)

	// Resolve Unix address
	unixAddress, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		_ = catcher.Error("cannot resolve unix address", err, nil)
		os.Exit(1)
	}

	// Listen on Unix socket
	listener, err := net.ListenUnix("unix", unixAddress)
	if err != nil {
		_ = catcher.Error("cannot listen to unix socket", err, nil)
		os.Exit(1)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	plugins.RegisterCorrelationServer(grpcServer, &correlationServer{})

	// Serve gRPC
	if err := grpcServer.Serve(listener); err != nil {
		_ = catcher.Error("cannot serve grpc", err, nil)
		os.Exit(1)
	}
}

func (c *correlationServer) Correlate(ctx context.Context, alert *plugins.Alert) (*emptypb.Empty, error) {
	// Implement your correlation logic here
	// Example: Correlate multiple failed logins from same IP to create incident

	return nil, nil
}
```

### Developing a Notification Plugin

Notification plugins send notifications to UTMStack and external systems (e.g., Slack channels, ticketing systems).

```go
package main

import (
	"context"
	"net"
	"os"

	"github.com/threatwinds/go-sdk/catcher"
	"github.com/threatwinds/go-sdk/plugins"
	"github.com/threatwinds/go-sdk/utils"
	"google.golang.org/grpc"
)

type notificationServer struct {
	plugins.UnimplementedNotificationServer
}

func main() {
	// Create socket directory
	filePath, err := utils.MkdirJoin(plugins.WorkDir, "sockets")
	if err != nil {
		_ = catcher.Error("cannot create socket directory", err, nil)
		os.Exit(1)
	}

	// Create socket path for notification plugin
	socketPath := filePath.FileJoin("my_notification.sock")
	_ = os.Remove(socketPath)

	// Resolve Unix address
	unixAddress, err := net.ResolveUnixAddr("unix", socketPath)
	if err != nil {
		_ = catcher.Error("cannot resolve unix address", err, nil)
		os.Exit(1)
	}

	// Listen on Unix socket
	listener, err := net.ListenUnix("unix", unixAddress)
	if err != nil {
		_ = catcher.Error("cannot listen to unix socket", err, nil)
		os.Exit(1)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	plugins.RegisterNotificationServer(grpcServer, &notificationServer{})

	// Serve gRPC
	if err := grpcServer.Serve(listener); err != nil {
		_ = catcher.Error("cannot serve grpc", err, nil)
		os.Exit(1)
	}
}

func (n *notificationServer) Notify(ctx context.Context, msg *plugins.Message) (*emptypb.Empty, error) {
	// Implement your notification logic here
	// Example: Send alert to Slack channel

	return &emptypb.Empty{}, nil
}
```

### Plugin Configuration

Plugins can be configured using environment variables or configuration files. The go-sdk provides utilities for reading
configuration values:

```go
value := plugins.PluginCfg("plugin_name", false).Get("my_config_key").String()
```

### Plugin Testing

It's important to test your plugins thoroughly before deploying them to production. Here are some testing strategies:

1. **Unit Testing**: Write unit tests for your plugin logic.
2. **Integration Testing**: Test your plugin with the EventProcessor in a development environment.
3. **Load Testing**: Test your plugin under load to ensure it can handle the expected event volume.

## Scaling the System

This section provides guidance on scaling the EventProcessor and UTMStack ecosystem to handle increasing event volumes.

UTMStack can scale both vertically and horizontally to accommodate growing workloads and data volumes.

### Horizontal Scaling

For horizontal scaling, UTMStack only requires adding more nodes to Docker Swarm:

1. **Add Docker Swarm Nodes**: Simply add more nodes to your Docker Swarm cluster.
2. **OpenSearch Node Affinity**: Add more OpenSearch nodes with node affinity to ensure proper data distribution.
3. **Automatic Component Scaling**: The rest of the components will automatically scale based on the number of Docker
   nodes available.

This approach simplifies scaling operations and ensures that your UTMStack deployment can grow seamlessly with your
needs.

### Vertical Scaling

You can also scale the system vertically by increasing the resources of existing nodes:

1. **Increase CPU**: Add more CPU cores to handle more events.
2. **Increase Memory**: Add more memory to improve caching and reduce disk I/O.
3. **Faster Disks**: Use faster disks (e.g., SSDs) to improve I/O performance.

### Database Scaling

The OpenSearch database is a critical component that can be scaled to handle more data:

1. **Add OpenSearch Nodes**: Add more OpenSearch nodes to the cluster with proper node affinity settings.
2. **Increase Shards**: Increase the number of shards to distribute the data more effectively.

### Queue Scaling

UTMStack handles event volume spikes through its built-in scaling mechanisms without relying on external message queue
systems like Kafka, RabbitMQ, or Redis. The system's architecture is designed to efficiently process events directly
within the Docker Swarm infrastructure.

### Monitoring and Alerting

Implement monitoring and alerting to detect scaling issues:

1. **Prometheus**: Use Prometheus to collect metrics.
2. **Grafana**: Use Grafana to visualize metrics and set up alerts.

## Migration and Upgrades

Migration and upgrades in UTMStack are managed automatically by the UTMStack installer, eliminating the need for manual
intervention. The installer handles all aspects of the upgrade process, including:

1. **Version Compatibility Checks**: The installer automatically verifies compatibility between components.
2. **Database Schema Updates**: Any required schema changes are applied automatically.
3. **Component Upgrades**: All components are upgraded in the correct order.

In future releases, UTMStack will include an updates service that will further streamline this process by:

1. **Automatic Update Detection**: Notifying administrators when updates are available.
2. **Scheduled Updates**: Allowing updates to be scheduled during maintenance windows.
3. **Rolling Updates**: Implementing updates with minimal downtime.
4. **Update Verification**: Verifying the success of updates and automatically rolling back if issues are detected.

This automated approach ensures that your UTMStack deployment remains up-to-date with minimal effort and risk.

## Community Resources and Support

This section provides information on community resources and support options for the EventProcessor and UTMStack
ecosystem.

### Documentation

Access comprehensive documentation:

1. **Official Documentation**: Visit the official documentation
   at [UTMStack Documentation](https://documentation.utmstack.com).
2. **GitHub Repositories**: Check the GitHub repositories for READMEs and wikis.
3. **Code Comments**: Look at the code comments for detailed information about specific functions.

### Community Forums

Engage with the community:

1. **GitHub Discussions**: Participate in discussions on GitHub and get community support.
2. **Discord**: Join the UTMStack Discord server for development discussions and collaboration.

### Training and Tutorials

Learn from training materials and tutorials:

1. **Official Tutorials**: Follow the official tutorials on the UTMStack website.
2. **YouTube Channel**: Watch tutorial videos on the UTMStack YouTube channel.
3. **Webinars**: Attend webinars on new features and best practices.

### Commercial Support

Get commercial support if needed:

1. **Professional Services**: Engage professional services for implementation and customization.
2. **Technical Support**: Purchase technical support for production environments.
3. **Training Services**: Get training for your team from UTMStack experts.

### Contributing

Contribute to the project:

1. **Bug Reports**: Report bugs on GitHub.
2. **Feature Requests**: Submit feature requests on GitHub.
3. **Pull Requests**: Contribute code through pull requests.
4. **Documentation**: Help improve the documentation.
5. **Community Support**: Help answer questions from other users.
