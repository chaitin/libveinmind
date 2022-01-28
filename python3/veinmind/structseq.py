def describe(structseq, fields):
	"Retrieve dict from fields to indices of structseq."

	result = dict()
	stub = structseq([i for i in range(structseq.n_fields)])
	for field in fields:
		if not hasattr(stub, field):
			continue
		result[field] = getattr(stub, field)
	return result
