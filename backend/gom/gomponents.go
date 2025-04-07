package gom

import (
	"io"
	"strconv"
	"todo/records"

	. "maragu.dev/gomponents"      // Import Gomponents
	. "maragu.dev/gomponents/html" // Import Gomponents
)

func Completed(w io.Writer, rh *records.RecordHandler) {
	records := rh.GetCompletedRecords()
	listItems := make([]Node, len(records))
	for i, record := range records {
		listItems[i] = Li(
			Class("flex items-center justify-between p-2 border-b"),
			Text(record.Task),
		)
	}

	page := HTML(
		Head(
			Title("Simple Gomponents HTML Page"),
			Script(
				Src("https://cdn.tailwindcss.com"),
			),
		),
		Body(
			Class("bg-gray-100 min-h-screen"),
			A(
				Class("block text-center mt-4"),
				Attr("href", "/"),
				Button(
					Class("px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"),
					Attr("type", "button"),
					Text("Back to Home"),
				),
			),
			Div(
				Class("max-w-2xl mx-auto mt-8 p-6 bg-white rounded-lg shadow-lg"),
				H1(
					Class("text-2xl font-bold mb-4 text-gray-800"),
					Text("Completed Tasks"),
				),
				Ul(
					append([]Node{Class("space-y-2")}, listItems...)...,
				),
			)),
	)

	page.Render(w)
}

func Home(w io.Writer, rh *records.RecordHandler) {
	records := rh.GetActiveRecords()
	listItems := make([]Node, len(records))
	for i, record := range records {
		listItems[i] = Li(
			Class("flex items-center justify-between p-2 border-b"),
			Text(record.Task),
			Button(
				Class("px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"),
				Attr("type", "button"),
				Attr("onclick", "fetch('/complete/"+strconv.Itoa(record.RecordID)+"', { method: 'GET' }).then((res)=>location.reload());"),
				Text("Complete"),
			),
		)
	}

	page := HTML(
		Head(
			Title("Simple Gomponents HTML Page"),
			Script(
				Src("https://cdn.tailwindcss.com"),
			),
		),
		Body(
			Class("bg-gray-100 min-h-screen"),
			A(
				Class("block text-center mt-4"),
				Attr("href", "/completed"),
				Button(
					Class("px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"),
					Attr("type", "button"),
					Text("View Completed Tasks"),
				),
			),
			Div(
				Class("max-w-2xl mx-auto mt-8 p-6 bg-white rounded-lg shadow-lg"),
				Form(
					Action("/add"),
					Method("POST"),
					Attr("id", "taskForm"),
					Class("mb-6"),
					Div(
						Class("flex gap-2"),
						Input(
							Class("flex-1 px-4 py-2 border rounded focus:outline-none focus:ring-2 focus:ring-blue-500"),
							Attr("type", "text"),
							Attr("id", "nameInput"),
							Attr("placeholder", "Enter task details"),
						),
						Button(
							Class("px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"),
							Attr("type", "submit"),
							Text("Add Task"),
						),
					),
				),
				Ul(
					append([]Node{Class("space-y-2")}, listItems...)...,
				),
			)),
		Script(
			Raw(` document.getElementById('taskForm').onsubmit = (event)=> { event.preventDefault(); var task = document.getElementById('nameInput').value; 
			fetch('/add', { method: 'POST', body: task, headers: { 'Content-Type': 'text/plain' } }).then((res)=>location.reload()); };`),
		),
	)

	page.Render(w)
}
