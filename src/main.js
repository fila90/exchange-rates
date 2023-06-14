import Chart from "chart.js/auto";
import { format, parseISO } from "date-fns";

let myChart;

fetch("/api/exchange")
  .then((resp) => resp.json())
  .then((resp) => {
    console.log(resp);
    generateChart(resp);
  })
  .catch((err) => console.error(err));

function generateChart(data) {
  if (myChart) {
    myChart.destroy();
  }

  myChart = new Chart(document.querySelector("#chart"), {
    type: "line",
    data: {
      labels: data.erste.map((rate) =>
        format(parseISO(rate.date), "dd/MM/yyyy")
      ),
      datasets: [
        {
          backgroundColor: "rgba(40,112,237,.7)",
          borderColor: "rgba(40,112,237,.7)",
          data: (data.erste ?? []).map((rate) => rate.rate),
          label: "Erste",
          tension: 0.1,
        },
        {
          backgroundColor: "#e2001a",
          borderColor: "#e2001a",
          data: (data.ucb ?? []).map((rate) => rate.rate),
          label: "UCB",
          tension: 0.1,
        },
      ],
    },
  });
}
