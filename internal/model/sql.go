package model

const (
	GetNextAssignee = `
with a as(
    select support.id, case when b.count_el is null then 0 else b.count_el end as count_el
    from support
             left join (select assignee_id, count(*) as count_el
                        From incident
                        Group by assignee_id) b on support.id=b.assignee_id
)
select id
from a
where count_el = (select min(count_el) from a)
`
)
