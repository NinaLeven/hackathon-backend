package model

const (
	GetNextAssignee = `
with a as (
    select assignee_id, count(*) as count_el from incident
    group by assignee_id
)
select assignee_id from a
where count_el = (select min(count_el) from a)
`
)
