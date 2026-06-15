## Notes
- When you need to create new error you should use apperrors.New function.
- All times when you check "if err != nil" you should wrap existing err to new err using apperrors.Wrap with additional message.
- All error messages start with a lowercase letter
