import Chart from "chart.js/auto";
import { format, parseISO } from "date-fns";

let myChart;
const rateErste = document.querySelector(".rate--erste");
const rateUcb = document.querySelector(".rate--ucb");

fetch("/api/exchange")
  .then((resp) => resp.json())
  .then((resp) => generateChart(resp))
  .catch((err) => console.error(err));

function generateChart(data) {
  const { erste = [], ucb = [] } = data;

  rateErste.innerText = `Erste ~ ${erste.at(-1).rate}`;
  rateUcb.innerText = `UCB ~ ${ucb.at(-1).rate}`;

  if (myChart) {
    myChart.destroy();
  }

  myChart = new Chart(document.querySelector("#chart"), {
    type: "line",
    data: {
      labels: data.erste.map((rate) => format(parseISO(rate.date), "dd/MM")),
      datasets: [
        {
          backgroundColor: "rgba(40,112,237,.7)",
          borderColor: "rgba(40,112,237,.7)",
          data: erste.map((rate) => rate.rate),
          label: "Erste",
        },
        {
          backgroundColor: "#e2001a",
          borderColor: "#e2001a",
          data: ucb.map((rate) => rate.rate),
          label: "UCB",
        },
      ],
    },
  });
}
