package aegis.mcp

import rego.v1

default allow = false
default deny = false
default final_allow = false

# Allow if any rule matches
allow if {
    some rule in data.rules
    rule.effect == "allow"
    matches(rule)
}

# Deny if any rule matches and effect is deny (Deny overrides Allow)
deny if {
    some rule in data.rules
    rule.effect == "deny"
    matches(rule)
}

matches(rule) if {
    # Check method
    rule.methods[_] == input.method

    # Check tool if applicable
    check_tool(rule)

    # Check roles
    check_roles(rule)

    # Check agent IDs
    check_agent_ids(rule)

    # Check semantic safety
    input.inspection.safety_score >= data.semantic.minimum_safety_score
    check_intents(rule)
}

check_tool(rule) if {
    not rule.tools
}

check_tool(rule) if {
    rule.tools[_] == input.tool
}

check_roles(rule) if {
    not rule.roles
}

check_roles(rule) if {
    rule.roles[_] == input.auth.roles[_]
}

check_agent_ids(rule) if {
    not rule.agent_ids
}

check_agent_ids(rule) if {
    rule.agent_ids[_] == input.auth.agent_id
}

check_intents(rule) if {
    not rule.intent_allow_list
}

check_intents(rule) if {
    count({intent | intent := input.inspection.intent_categories[_]; not contains_intent(rule.intent_allow_list, intent)}) == 0
}

contains_intent(list, item) if {
    list[_] == item
}

final_allow if {
    allow
    not deny
}

# Final Decision
decision := {
    "allowed": allow,
    "denied": deny,
    "final_allow": final_allow
}
