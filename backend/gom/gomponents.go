package gom

import (
	"io"
	"strconv"
	"todo/records"

	. "maragu.dev/gomponents"
	. "maragu.dev/gomponents/html"
)

func Home(w io.Writer, rh *records.RecordHandler) {
	activeRecords := rh.GetActiveRecords()
	completedRecords := rh.GetCompletedRecords()

	activeItems := make([]Node, len(activeRecords))
	for i, record := range activeRecords {
		activeItems[i] = Li(
			Class("flex items-center justify-between p-3 bg-gray-50 rounded-lg border border-gray-200"),
			Span(Class("text-gray-800 flex-1"), Text(record.Task)),
			Button(
				Class("ml-3 px-3 py-1.5 bg-green-500 text-white text-sm rounded-lg hover:bg-green-600 transition-colors"),
				Attr("type", "button"),
				Attr("onclick", "fetch('/complete/"+strconv.Itoa(record.RecordID)+"',{method:'GET'}).then(()=>location.reload())"),
				Text("✓ Complete"),
			),
		)
	}

	completedItems := make([]Node, len(completedRecords))
	for i, record := range completedRecords {
		completedItems[i] = Li(
			Class("flex items-center p-3 bg-green-50 rounded-lg border border-green-100"),
			Span(Class("text-green-500 mr-3 font-bold"), Text("✓")),
			Span(Class("text-gray-500 line-through"), Text(record.Task)),
		)
	}

	var activeContent Node
	if len(activeRecords) == 0 {
		activeContent = P(Class("text-center text-gray-400 py-10"), Text("No active tasks — add one above!"))
	} else {
		activeContent = Ul(append([]Node{Class("space-y-2")}, activeItems...)...)
	}

	var completedContent Node
	if len(completedRecords) == 0 {
		completedContent = P(Class("text-center text-gray-400 py-10"), Text("No completed tasks yet."))
	} else {
		completedContent = Ul(append([]Node{Class("space-y-2")}, completedItems...)...)
	}

	activeCount := strconv.Itoa(len(activeRecords))
	completedCount := strconv.Itoa(len(completedRecords))

	page := HTML(
		Head(
			Title("Todo App"),
			Script(Src("https://cdn.tailwindcss.com")),
		),
		Body(
			Class("bg-gradient-to-br from-blue-50 to-indigo-100 min-h-screen"),
			Div(
				Class("max-w-2xl mx-auto pt-10 px-4 pb-8"),
				H1(
					Class("text-3xl font-bold text-center text-indigo-800 mb-8"),
					Text("My Tasks"),
				),
				Div(
					Class("bg-white rounded-2xl shadow-xl overflow-hidden"),
					Div(
						Class("flex border-b border-gray-200"),
						Button(
							Attr("id", "tabActive"),
							Attr("type", "button"),
							Attr("onclick", "switchTab('active')"),
							Class("flex-1 py-4 px-6 text-sm font-semibold text-indigo-600 border-b-2 border-indigo-600 bg-indigo-50 transition-colors"),
							Text("Active"),
							Span(Class("ml-2 px-2 py-0.5 bg-indigo-600 text-white text-xs rounded-full"), Text(activeCount)),
						),
						Button(
							Attr("id", "tabCompleted"),
							Attr("type", "button"),
							Attr("onclick", "switchTab('completed')"),
							Class("flex-1 py-4 px-6 text-sm font-semibold text-gray-500 hover:text-gray-700 transition-colors"),
							Text("Completed"),
							Span(Class("ml-2 px-2 py-0.5 bg-gray-200 text-gray-600 text-xs rounded-full"), Text(completedCount)),
						),
					),
					Div(
						Attr("id", "panelActive"),
						Class("p-6"),
						Form(
							Attr("id", "taskForm"),
							Class("mb-5"),
							Div(
								Class("flex gap-2"),
								Input(
									Class("flex-1 px-4 py-2.5 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"),
									Attr("type", "text"),
									Attr("id", "nameInput"),
									Attr("placeholder", "What needs to be done?"),
									Attr("autocomplete", "off"),
								),
								Button(
									Class("px-5 py-2.5 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors font-medium"),
									Attr("type", "submit"),
									Text("Add"),
								),
							),
						),
						activeContent,
					),
					Div(
						Attr("id", "panelCompleted"),
						Class("p-6 hidden"),
						completedContent,
					),
				),
			),
		),
		Script(Raw(`
document.getElementById('taskForm').onsubmit = e => {
	e.preventDefault();
	const task = document.getElementById('nameInput').value.trim();
	if (!task) return;
	fetch('/add', { method: 'POST', body: task, headers: { 'Content-Type': 'text/plain' } })
		.then(() => location.reload());
};

function switchTab(tab) {
	const isActive = tab === 'active';
	document.getElementById('panelActive').classList.toggle('hidden', !isActive);
	document.getElementById('panelCompleted').classList.toggle('hidden', isActive);
	const on  = 'flex-1 py-4 px-6 text-sm font-semibold text-indigo-600 border-b-2 border-indigo-600 bg-indigo-50 transition-colors';
	const off = 'flex-1 py-4 px-6 text-sm font-semibold text-gray-500 hover:text-gray-700 transition-colors';
	document.getElementById('tabActive').className    = isActive ? on : off;
	document.getElementById('tabCompleted').className = isActive ? off : on;
	const url = new URL(window.location);
	isActive ? url.searchParams.delete('tab') : url.searchParams.set('tab', 'completed');
	history.replaceState(null, '', url);
}

if (new URLSearchParams(window.location.search).get('tab') === 'completed') switchTab('completed');
`)),
	)

	page.Render(w)
}
