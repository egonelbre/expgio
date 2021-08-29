package main

func NewDemoDiagram() *Diagram {
	diagram := NewDiagram()

	issue := diagram.NewNode(Label("issue"), V(1, 1), V(8, 1),
		nil,
		[]*Port{{Name: "done"}},
	)
	schedule := diagram.NewNode(Label("schedule"), V(1, 3), V(8, 1),
		nil,
		[]*Port{{Name: "done"}},
	)
	issueComment := diagram.NewNode(Label("issue-comment"), V(1, 5), V(8, 1),
		nil,
		[]*Port{{Name: "done"}},
	)
	_ = diagram.NewNode(Label("noop"), V(1, 7), V(8, 1),
		nil,
		[]*Port{{Name: "done"}},
	)

	manageLabels := diagram.NewNode(Label("manage-labels"), V(12, 1), V(8, 1),
		[]*Port{{Name: "start"}},
		[]*Port{{Name: "done"}},
	)
	composeComment := diagram.NewNode(Label("compose-comment"), V(12, 3), V(8, 1),
		[]*Port{{Name: "start"}},
		[]*Port{{Name: "done"}},
	)
	issueFlow := diagram.NewNode(List{"edited", "closed", "deleted"}, V(12, 5), V(8, 3),
		[]*Port{{Name: "start"}},
		[]*Port{{Name: "done"}},
	)
	daily := diagram.NewNode(Label("daily"), V(12, 9), V(8, 1),
		[]*Port{{Name: "start"}},
		[]*Port{{Name: "done"}},
	)

	updateLabels := diagram.NewNode(List{"add-labels", "remove-labels"}, V(22, 1), V(8, 2),
		[]*Port{{Name: "start"}},
		[]*Port{{Name: "done"}},
	)
	createComment := diagram.NewNode(Label("create-comment"), V(22, 4), V(8, 1),
		[]*Port{{Name: "start"}},
		[]*Port{{Name: "done"}},
	)
	stale := diagram.NewNode(Label("stale"), V(22, 6), V(8, 2),
		[]*Port{{Name: "start"}, {Name: "interval"}},
		[]*Port{{Name: "done"}, {Name: "fail"}},
	)

	diagram.Conns = []*Conn{
		{From: issue.Out[0], To: manageLabels.In[0]},
		{From: issue.Out[0], To: composeComment.In[0]},
		{From: issue.Out[0], To: issueFlow.In[0]},

		{From: schedule.Out[0], To: daily.In[0]},
		{From: issueComment.Out[0], To: stale.In[0]},

		{From: manageLabels.Out[0], To: updateLabels.In[0]},
		{From: manageLabels.Out[0], To: createComment.In[0]},

		{From: composeComment.Out[0], To: createComment.In[0]},

		{From: daily.Out[0], To: stale.In[0]},
	}

	return diagram
}
