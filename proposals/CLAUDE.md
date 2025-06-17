# CLAUDE.md - Proposals Directory

This file provides guidance to Claude Code when working with proposal documents in this directory.

## Directory Structure

```
proposals/
├── prds/                     # Product Requirements Documents
├── design-docs/              # Technical Design Documents  
├── templates/                # Document templates
└── CLAUDE.md                # This file
```

## Document Types and Purpose

### Product Requirements Documents (PRDs)
**Location**: `proposals/prds/`
**Purpose**: Define WHAT to build and WHY
**Focus**: User needs, business value, success criteria
**Audience**: Product managers, stakeholders, business users

### Design Documents
**Location**: `proposals/design-docs/`
**Purpose**: Define HOW to build something
**Focus**: Technical implementation, architecture, components
**Audience**: Engineers, architects, technical reviewers

### Templates
**Location**: `proposals/templates/`
**Purpose**: Standardized formats for creating consistent proposals
**Usage**: Copy template, rename, and fill in content

## Naming Convention

Use the format: `YYYY-MM-feature-name.md`

Examples:
- `2025-06-initial-streaming-api.md`
- `2025-07-section-subscriptions.md`
- `2025-08-mutation-interface.md`

## Document Relationship

PRDs and Design Docs should be paired and cross-referenced in frontmatter:

**PRD frontmatter:**
```yaml
---
title: "Feature Name Product Requirements"
author: ["Author Name"]
date: "2025-06-17"
status: "Draft"
design_doc: "proposals/design-docs/2025-06-feature-name.md"
---
```

**Design Doc frontmatter:**
```yaml
---
title: "Feature Name Design Document"  
author: ["Author Name"]
date: "2025-06-17"
status: "Draft"
prd: "proposals/prds/2025-06-feature-name.md"
---
```

## Status Values

Use these standard status values:
- **Draft**: Initial version, work in progress
- **Review**: Ready for stakeholder feedback
- **Approved**: Finalized and approved for implementation
- **Implemented**: Feature has been built and deployed

## When to Create Proposals

### Always Create Both PRD + Design Doc For:
- New major features or capabilities
- Breaking changes to existing APIs
- Cross-team initiatives
- Features that impact user experience
- Architectural changes

### PRD Only For:
- Minor feature enhancements with obvious implementation
- Configuration changes
- Simple bug fixes that affect user experience

### Design Doc Only For:
- Pure technical refactoring
- Performance optimizations
- Internal architecture improvements
- Developer tooling enhancements

## Writing Guidelines

### For PRDs:
- **Start with user value** - why does this matter to users?
- **Be specific about success** - include measurable criteria
- **Focus on outcomes** - what users will be able to accomplish
- **Include concrete examples** - real user scenarios
- **Consider all stakeholders** - not just end users

### For Design Docs:
- **Technical depth** - enough detail for implementation
- **Decision rationale** - explain why you chose this approach  
- **Consider alternatives** - show you evaluated multiple options
- **Plan for testing** - think about verification upfront
- **Implementation ready** - actionable tasks and phases

### For Both:
- **Use templates** - start with the provided templates
- **Clear headings** - make documents scannable
- **Link related content** - reference other docs, code, issues
- **Include examples** - show concrete usage scenarios
- **Cross-reference** - link PRDs and design docs together

## Review Process

1. **Draft**: Author creates initial version using templates
2. **Internal Review**: Team reviews for completeness and clarity
3. **Stakeholder Review**: Broader review by affected teams/users
4. **Approval**: Document approved and ready for implementation
5. **Implementation**: Work proceeds according to the plan
6. **Update Status**: Mark as "Implemented" when complete

## Common Patterns

### Feature Development Flow:
1. Identify user need or technical requirement
2. Create PRD defining the user value and requirements
3. Create Design Doc defining the technical approach
4. Review both documents with stakeholders
5. Implement according to the design
6. Update documents if significant changes occur

### Architecture Decision Flow:
1. Create Design Doc exploring alternatives
2. Include decision rationale and trade-offs
3. Review with architecture team
4. Document final decision and approach
5. Reference in implementation

## Tips for Claude Code

### When Analyzing Existing Proposals:
- Read both PRD and Design Doc to understand full context
- Check status to understand implementation state
- Look at cross-references to understand relationships
- Consider how new work relates to existing proposals

### When Creating New Proposals:
- Always start with templates from `proposals/templates/`
- Use consistent naming convention
- Include proper frontmatter with cross-references
- Follow the established patterns in existing documents
- Reference related work and existing architecture

### When Implementing Features:
- Refer to both PRD and Design Doc for complete understanding
- Update document status as work progresses
- Note any significant deviations from the original design
- Keep implementation aligned with the documented approach

### When Making Changes:
- Update related proposal documents if scope changes
- Maintain consistency between PRD and Design Doc
- Document architectural decisions that affect future work
- Consider impact on related proposals and features

## Quality Standards

### All Documents Should:
- Have complete frontmatter with proper cross-references
- Follow the template structure and include all sections
- Be written clearly with concrete examples
- Include specific, measurable success criteria
- Link to related documents and external references

### PRDs Should:
- Focus on user value and business outcomes
- Include specific use cases and user workflows
- Define clear success metrics and launch criteria
- Consider future enhancements and evolution

### Design Docs Should:
- Include sufficient technical detail for implementation
- Document key architectural decisions and alternatives
- Provide clear component interfaces and interactions
- Include comprehensive test plans and success criteria

Remember: These documents are the foundation for feature development. Take time to create thorough, well-structured proposals that will guide implementation and serve as reference for future work.