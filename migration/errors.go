package migration

import "fmt"

// LoopError - loop error
type LoopError struct {
	LastMigration string
}

func (e *LoopError) Error() string {
	return fmt.Sprintf("Loop detected: last migration in the loop '%v'", e.LastMigration)
}

// RequirementNotFoundError - requirement error
type RequirementNotFoundError struct {
	Migration   *Meta
	Requirement string
}

func (e *RequirementNotFoundError) Error() string {
	return fmt.Sprintf("Can not find requirement with name %v for %v", e.Requirement, e.Migration)
}

// RequirementDuplicateError - Migration meta requirement
type RequirementDuplicateError struct {
	Migration   *Meta
	Requirement *Meta
}

func (e *RequirementDuplicateError) Error() string {
	return fmt.Sprintf("Requirement duplicate in %v (%v)", e.Migration, e.Requirement)
}

// NullRequirement - requirement error
type NullRequirement struct {
	Migration *Meta
}

func (e *NullRequirement) Error() string {
	return fmt.Sprintf("Can not add null requirement to %v", e.Migration)
}

// NotFound - migration not found
type NotFound struct {
	Name   string
	Filter string
}

func (e *NotFound) Error() string {
	return fmt.Sprintf("Migration %v='%v' not found", e.Filter, e.Name)
}

// NoDBInformation - no db information for migration
type NoDBInformation struct {
	Migration *Meta
}

func (e *NoDBInformation) Error() string {
	return fmt.Sprintf("No information from database for migration '%v'", e.Migration.Name)
}

// NoFileInformation - no file information for migration
type NoFileInformation struct {
	Migration *Meta
}

func (e *NoFileInformation) Error() string {
	return fmt.Sprintf("No information from files for migration '%v'", e.Migration.Name)
}
