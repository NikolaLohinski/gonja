Persons:
{{ persons|map(attribute='FirstName')|join(', ') }}

Names:
{{ persons|map(attribute='LastName')|map('upper')|join(', ') }}
