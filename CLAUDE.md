## Notes
- When you need to create new error you should use apperrors.New function.
- All times when you check "if err != nil" you should wrap existing err to new err using apperrors.Wrap with additional message.
- All error messages start with a lowercase letter
- All time variables should have time.Time type, don't use int64 and Unix timestamp inside
- Order of klines in array is always: 0 is last kline; len()-1 is first kline