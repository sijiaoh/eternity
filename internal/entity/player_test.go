package entity

// Compile-time check: Player must implement Entity interface.
var _ Entity = (*Player)(nil)
