package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Get the secret key from user input.
	var secret string
	fmt.Printf("Please enter your secret key: ")
	fmt.Scanln(&secret)

	// Get the keypair of the quest account from the secret key.
	questAccount, _ := keypair.Parse(secret)

	// Fund and create the quest account.
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + questAccount.Address())
	if err != nil {
		log.Fatal(err)
	}

	// Get and print the response from friendbot.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Friendbot response:")
	fmt.Println(string(body))

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: questAccount.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatalln(err)
	}

	// Get the base64 encoding of the image.
	imgBase64 := "iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAYAAACqaXHeAAAgAElEQVR4nG2beYxlV17fP+fec9e31l7V1aurF3e393U8g2MYw2Q0M4JkgBAlEn+FCEWJkigRShRppCxS8kcWRREoBAkk+AMCSQgwNmQG8DCLPV7Guz12u9u9d1V1Vb1X9ba7nxud89ZqKKm7Xr1737nn/Jbv7/tbnvjVr38cUZYSAZQCQUmJfg3mvdFPKUZ/jt83N43uLkEIKEfXRClG9+urYnL7oSXL6drmzdFN42eL0U3leOF7rpdiZkE1WkaU5tZy5hHD26ZvluWhZ+XSEkIKkOaqJWa2OLPD2Z2Le68OD6t/W6KcuceiRJnjz8hx5jOlEdrkEXoRi+lhJ48VI+nqP8rJZyfPKc2jhgIvJ3qZHnhGyJPDj5WMQNpirFpxeNHZ04rykGD0QpZ1WInl6MmTMxltCISYUdUhE5gRtllwvLo1tcJZZZj1rNHm7/3c9Fn6lqEMy6nFiIk9HVZaCdoC7jnFVCDl7HuMFSFG65XT486c75CexYxL6BfWjH2KsbanGh7+Gi5mjV7NCkKMLaC8xwKEOGzujKyinJH32G1G5xKji9ISU1NTKCxLDM1xrKOx9EYmKMrp6crZXQgxulYONTHrQaMzGoyxZt4fC6mcMQghprJhuuashY4tYKxdfZgx5mDOMdHzSE3l5HMjjxmtUSL1hsTo6Ra2eWWNHzSjLTE+iTWz+8kBZ9xnLIiJccyYrDjsYVpQE62roRMbg1D3rDk+zkTA47WtGRuZWqMthgopx5KY+VEjjx5bg5xAwNhhpmseAqnxmSdam3XPe59izVwsxSGZjQxh8scQuMRU6uZ6OY04k81MdDnUbDn+i4lPH1pbv2vNPHfmeLNyNS4wflg53sTEzkfOUU4PfggOxYwJc/ieyWcnH56G2bF0xdh1GNp8ORtpJg+atZmxLd0D2KKcWPHQHWYWKkfGPnqOGANnObS8KQhOgO2viHliCl7W+AiH8GDmQJOgMZbECHHHUW7ks7N8YSy5Kc+YcZt7EXZ0TYzt+FCEmpHb+JoYouEYw4ZmMVaRdgGmhxAzoUrMBg3joiP0L6fSLsUsAxoB5kjCkwfOhkZr9JY1+xwmmp++ngrzEAZOZDLEgdnYX473NTHJGSsen2dG+OMlpWWVE9RlDEoz/jlcwDJMaYoBY0+cYWmH/HEET6MdTrigmFrPrJ3Phioxq/1yTITKqTjMdsTowNMfawooQzeYcZ1x6B4KYbJzA4iHecB4sdEDxYjRzZJKMXlVHuJPkzPMuMbU35n4ozUOkeJQEB0hVDklMaO1mASyWSwpJ1rnnn2MQ94k0E3WO0yt9I+teYAwseJeIZTTWGnMrRwh6DRPGJvrVH/ljD5nCM9sPBVD7UxlPorU43vVvVRqCuuz2FQeUsR0z2J2FxNwnLlnsvnRnZoJjk3YOsxTRx+YQd97TM5sX82C0vj+8rAHz2YfYwyfKFNNZfVXnH1y3JHgpi5TMrvqrBlNtP2XXHlqdWU5dSVLjNF68q+c/hsfvGRqTuU9O5zw7hGJYbqRCblj5ppB5HKCzOM1xzAyFnRrd5dLP/ju8PNqtAeNr6Pr1pgQiamJ/yUFifHeyxnOMMt7tKXr//QutDKUGiLDaGcTEBz/PT7oKJ6KcnrQGQEfFuasMY2jxThnVVNBl+WMlamS7Y/f4urv/Q9OvXyZpU5/CmL3WKg12pslyhlBl0MBTQ7PJCsdK8QahWU5eeg4RI0xwayrpjR41sNmOHI5+sxMOD9sIjNCmLUea0KSxZD6zhxKW+UDn32ena2Ug9c/4ZlvvUL7ofP8+bMPE9btyZpjajxce4o9szxIv/SIaFYdtrpycoaJACfqmmhoJgE5dJip39mU5t/Q7NQMos9Yzj3uJCb4oF+rqSUN7X5oruNcwFik4sHiEp2vHOE3zrr8u6//Gq994+vYlvgr3QzjHsUQy8aaH4XR7774+7zw73+J+2rX8FU6dUvNBHUSYlkKr2gRiYWJ+AzUlAW2tkjLnphRrfMO88lHhGHAzfI0veB+KAsExei3Mq5UZMWIrSkjrFIVFErR2b3GsdUlesHGjDamm+bWy5TCYfflb1I9/RDq6An+5A9/hwO/w5fPrFMUI76ixl5UkqUZfucdntgocMN53rh7H0kxqiuUJWefep4/f/t7dF9/iYdWTnGt8Sg71Ib298rv/2r29kdb8guPLdLOq3T7Cb2DHtGgS2eQ8+gjG2zeuEuRZ7iOjSti3KIDlkS6Hq1uSVHkZFlBrVpDCMv8i9LExFnHsbEs22xkq52Q9fc5/7kvsa1WJuUATZi00ESZcWTnRQaFw/4g0yKl19unoWIip0Il9EmzFCUkpaXXdSiEpJuWhMUei1ULx3PYyo8ibB9huTiuy+b1y8ir71Lc/zk2P/x/nDv3CP7Gj9M66OcyD07wxKPLFPk20p9jOVRcurHLi7/3IkiHl/7se3zpR+/n8YeOs7i+zs7N2xSporm4xOZeQrHwJFSP4gpBNKKjapSU5CXEasoPqhuKo/1XyG0brGAES8qYrrI8RJHg+x677kN82O1z4uQZos4u1p/+dzaf+Rvcd3zVrNT59BVOhG3i+kUybwU3Sog++SOOnV2lfdDhwWCP3W7M6snTfPrB2zx8coH2/Hm6O6/TrdaIi4JK9wN2Lu8hN/JXaUUl1cChKvpU5pZpWrskSUyn1yFPApJon8XmBkWvy/LKEll/D7IOp5ZDagdv0N163ZiiIRbSwbIssqzEla5hd0lWEKe5MfiXr9zh6c8+hxA5SthGUINPv8V9cwVJWqLKhJPhp4TLCds3LrFYWcE9VifZeYlmsMQnt3sshz3mXI+7yTZp9T4qssbKhVMcORZS8QqKwmL56BEECfL4AkVRsHaiinXiYZ7sd1GFIM+3ObKRIr/9/jZXb+3xxc/cz9XLeyTJdT66scvluzs8liiWogN+u+hTXVjBkQHXb2+Tlxmu6/PFp47TKxvctdYoCoXneShVkufKaN9SgiIbCsaWBTK5y8LispYSIk2xVG6Y3sriAlvhWfr7OxwdfB8raZJ7DeaXOuzu7fJReRbZuYbdbnDhwfPI/Y+xHcVR9xbJ4EXjR5VA0jmw2O/0yZXFfhIT9QYGstN0QDXKEMLGtixtc/iuzevvbiLnQ4v5jTksCpp2TDjnsbsH3SjiAiGX513+2iPHeeTMApbncPHCIsLy8V1Ja/MO1zc3Gdgl/Vix2U5ZqVsEUhD6jqkvJWlGokqubPVw0jbzQUHVirhTnmJl41ET83cO+uxe+w7vfXSFz9xnUexcYyuuECUZr776XWwLLjzwGN2tmAMRYbXaPHhfk/1be9xtXeWBi2d448MDLp47ze2bGY7ns7Xfoxy0WKpa5FnB2Y0qnVbLYJJjW8wtNHjq3CLi13/pp7JmI5BrR9d57+PLxiev3+7x2ps7DPKCpbk6Z457nD5bp2ZLWr0UmZcMckWnl9LrKXKtbtsmGeSEnoX0XHxfGmlrs89yReC4pFnOdifhqYfW+eGn1/CCOYpU4DqSghysGFcELNUkn95qk8QFlivp9VPiQYbnO+RKYEtJrkoDykkGO60uX/nio5y9cJr+1ffppoI3r/bZmMtxvJBms0o3LUgzRdyPTOVHabDO81w+eOF+3r3W5uFqnacffpI8Tmg2O/zBn/0f9tKUrUHEfH2Vr65dMFpvzNWJ45SdzoD+YGBASWtKh7u8yMjyksVmSKsTDfFAWKaEHjr6vpxvX9onL1N8Kdje2eYb379BXuQ8fGaNZx8/YQ794MZ9nFo/QujaRJki0fihUrJMGYEmWU4URUirpMi1y1XJW7f5+C+u44ce/X7OUa9EpZKDTov29i5RmpEXQ5D2PQtbWLTiElnzKqiyS5lGVESV1JYmjG0d7FM9fQ7VadE66JHmBZVqlSjOae132N/vYVmSPM8RtkOcZhz0IrIsI88y02xKrYJSJUgp6YmCQZTRbW3zjcu7RlBx/STt3geEvsedrQ6/+fW3WK2FrC80qIUht/sDKpWKCZXadBv1KkWR0bRs8iwkzXPqlQDPcej2IwJfYmEZTWtG0mg2EHlmPhOlilwpwiCk0x/QqAb81kufIJ45ezJ7+Myq/Df/4Ef51d99hXarTTcTfHB1gOU4LM/VcVXM6jI0Gj65KrDylCgu+e5bLdKiHKYPalg4qTq29oZh+uz4zGmMWZvH9ixuXNmlPUipBA7Hzz9BP6jzzd/9TUI/YK0Z0k1y5iouNc/l2PF5WrsDNm/u0Y5iPNc2B8gVhK7EsoVxsYwSx7GoBz4Pn14mTwocR5qIoqtftSDU3R+UZRlFBKGv+T+OtPh4J81lnqZEUYcfXLrC0oJHPazx3pUD5iu+8RPXFuYBYeix2gj46EabhbrH8pzLi6/tQL3Os1/4KsdPnOXTa5dotz/B95c5ur5BUCTsvPMCRZpw7OIXOfHZDT69+T43btzCXj9NuX+Ln/2Ff8ra/DpZnvHbv/FfuH73YJJa1xzJfDPEPbrBvoJAFrTutrh881PSJJ6U8bQglht1nnlkiSjrkThNfvDuTQZpbg6q7+tGGYMkMeFa5xoPHVvm7KMPYv/kk0e/RqEszbs6nQFHV+bJBbhWTunCT3/paXr9AxqNALfok/uLzNcrFGnEtT2H8z/+Zc7f/zStbsz29g9xQw/bDbmv5tG+/B2Tlz3713+ajc98lXcvfQchbhInMXt3WwhScBQPXvwRjhy/yJNPP8d777zCoHtAqxdjWw6nNlaplDHV5iLHH1ykKJfp7t4lifuTpEcLYC7wuX9eUeYFjbUjuI5gruKR5jBIcqRl06x4nFioGtx4cKPOXCVR9q//i+e/9tjphvXg2XU2VjxOzLs88/BRzp1YNGj56MUNvvjEAk9vVDh9Yh63sc4XzoacXbf549db1FaWaG9ep5ozxA/L4uceWOUB9yqtXsHR0yc5d/YYL3zrJa7d+BhhZWAXNIMK1VCysLiAXfbZ3PyAs0cWObJxhrd+8BqDQZ92P0JYkmPHV5mvChxLP6PCjUuXSJP+xAJcKVltVvmRiws4gc/y6rIxeWFLE2FWFmoEjiT0XFaXqpxadI0VxLFScnWthqMqbHYtGivrrMx5SEti2/ucOrVEvd6EosN+K6Z+4SLO3ibVpSVEV5n43G3tEDjQLe7QyW2+vF5FczA/rGCFBX/2xg7/980/pd2ReH5JGFosy4DHHzzLieMXULaNLyskCgbRALKEhz77BB9+/3V2d/a4dO0OWZrzkz/xGOQFef/OMLli2lfQvr3QCPjbz62jMkXt3EV+4w/eNOFzbaXOUt1lLdAA6OI4sLvTxrZtLMdDKsshcAqWD/YpvRr9vQ5+vY4jPcKKpBflaMI2yANqtTWE3cGuL+IFDp77Ousbm+SFx60bDf7eZ1apCov9fsw3bxzwyvUYuwix7ZxYJIYV/sTTP8pXnv8cFInZhIZQXZwwGEqTjz+5RJwq7n/mEd783gcM9rb49M5d/vwv3uI//6Mv0jvpc194jv/4hx06g9QAneM4uLZNtVYn6UaIXPH4yTpPbVRMhNDA6tmSgVNB9A94J4s4cuYUxxd95GKjQl4I5JzUlAVneRWnuYTUcLu/aaS2Mn+EtL2LULEpj6s0RfhVNKmUZCxVFT/37DnyxKaXRHz7rsdmcYylk3N0uxZRliNlQZYl/OmbFh/efpXSTkxcdz1p6LOwhUl1o8jCap4laCoufv48n779Gq3Lb/DutW3++a98k6/9/HM8eu4I/+pnP8O//d3vG26hzdmrhNhz84SBonAcHj63aoo8KuohbBfpVdi83iYIaqwfW2P7xh1uXI6QdmONot9HzgUUWY5XqWNJm6zbwbUE0rFwixh3aZG+JSm071nzqNoKuXLw84d4xlcksWSQW/xJ+jjbgcKWDl0xBw2Bq0PVKCPc0hwjB5GPChkxBuzyOMOxJaVdEpw+o3khtSXBxdPPcOPlP+D6d/4Xb1+5zT/8ry/wa//yqzzuufz9L5X88h++aqzHdVwUDm69Oiy3SgfPE4hgDqs6Z8DxdKbo9SLCuQbNlVXad+5iFZrFuTWkX8P1PCy7pExivMaSyfdNycx1sYoYlath3y3eN6FtsebzuZW6SZLw5vh6+Rx3OILyj5FYK/jSI7BdPMvHE/q1T2Dr3w6+rV3Po+aG2PGAuqXwbYvA940gPFuwLGNzz8Xn/w6Pf/Ufk5cWN3Za/N1//TvsJuD6NR47f8FEgSSKSVtdLNfBDUP8wEPJKomokskqVn2R6vHjNNePUA9cwriDJxSyjA6w3QrW/Bpl5iOKzBxUWAVq0EHU6pPGotPdxtq/SdqaN6Wpn3p6jV57h4Ngie/wDLF0qdnTwYliVFfUJMmRkiJXSHuYHyhNnXsHuColvHOFO7cinv38OVg9RqRclpMWD3OF3xo8bjLLo489jxOEfP9//ie29vb5Z7/yIk9evGAyPM/1DN3OvUBLxdBv6ovYrsOVD67S7CUsry2RDQbkaWHScL85x/n1VWRSlrhxn2LvBsqSUEp67T5SdhgMcgNQWiiZTkJsj/UwI+/3sRs5nf0e+fIxvh38OLFyidPUILIjtT+X2NIy5q19XBMQ6dmjmp8udSvS7UvMqzvsF0s889STLNdr2L2Y2oLgtUGd9/OnsZ3hAFNRlqyde5yl+RVu3rzEdrvLC997jY31E+RFYULkwpnT2JaHihNKyzWp9rmNVYTrk/c6qKhPP4GlhrbslDI9wEpLm7LSoAybpGmKXa0yf2yN+vICwVxNOyjCryM16BUZjbpLsLxKZXmNn/mFH+OD4DEy4RqL8FyJY9vmwIUp3mHqBEqNkF4Ni6FaOBoV4uoy2+E6anWNxXNz1E4E1I5XeK8tUGlkOHyWF8RZYUKhjlhL559G2o4BPu0qO/t7nJuvc6TuoqVVjlo8o4ebqFCkBcLxzFtpt0856EI4B2FjOB8w6Pex+n1EkSOUQIQVst1t0H/rSo/lUOoCRpmZaq3BgSLjG28X3BErZiIjM+ZtjzsSOK6NHsBSYthrGHdw0pFg9BvVtRPY9n3m9cvbGe5ubu7LlY0UNmteRLPu62qKOawlcsJTi1yt1NnvtEzslNJmdxBRrS8Tbe+YMSTH1sNSXYRt8cabl3j82cdNGJT1Bgt+hX7ngIquK2IjfV0g8CvESYLjNbF14tDbx/a1KQ3Q3eO80EUcD2vpPpL33sTd2cYNXO788Cb2ypOURWH8LtOZoc7cpDTIXOhan2Ji9hoPNHnS1/RoinYLVWQGI/Q9mrYK8z70Bz1+8MIf8Wv/5CHWVyqGYQqh+G8vbw1becIyFah6WEE6Hrbng+UT7+9iLSxiSUE+6LG6voRIBygdYcoSRyj8hXmEH2C5HpKgYVJWx7JRaUJeKGNKOL4pbedJTNFPhlMvskW9UcVtLqGax6jM7Zhq7nhsTodPVQ4H3ZQpjY46MNZohG2kfFtMmyi2rhgLZbrGQzcpjSUVeUKn00ZmLYqDCMv1IemwdeuWAVBNohzHxXdcqr5PdS5A9Vrk3T3kiTOIvIfdXOFUJTU9QMMarWEWKdKYIrPJ0xypJaIBoyhiKHNTnLBcQZl0yCwXlcXkydAnnfmjuLUQpzmHCho0AomKc2xbkmaZ0bbWZqqbjpY1Hhwx/QCDJdo/1XSGyzKprG0SFWvUoyxGLXLV3SeOWnTbbbw8NkLL8pTt/YHBEtt28BwPTzqGc6zMuVQ3HsXZvmZKYZa0UHlmlFH0B4hKA5nHlHmKKgpEoQwQS9Xvgo7FUlLGGUiPFEm8exuZ6/JRFT+skdkO2aCNE0eo7SvQOELTs3HjIcJryjkeGCnKoSYdG2O6Ktflq+F4kiptkrRAymGtTs/o6vCoCy7jfo9miHtbW7QPdvnlP77E3/zsCeJMsXnzJr0oMsCq6a+hwK6DKgs2Fj3U3jWE6xl0L8N5A4JaKzqUkyaIShWVapJ3QODZlDrCDNKSstM2JiY83Ujw8MIGRTNF3ezRbbVZXvKHXRa/jtIInGZYyR6+pbBVjuf5RmsafHTZSf9kSlfElTmMDmGab+nmydjMtbZ19Vhjgo4KGsx0qc11HCxPMuhuc3Bwlzc/TKl7ivPr8yjbY6unTHPEdWsm4mjry/OCqgPJoE822MWrVNDZVV7YlBrUDaMNodJEFV1KT5Aoz4RlWTtygs7t23ihZ/wySRTJ7nWax48SrKU0rIS4FxuaaecZIsqQ68cR9XWshQFiKyJ3fFLNyUOPJFN40sLXwlAWllL4Wttq2CjRbhE4tjm8Zn55qut6BbqI5Wl30GmqFlocm6LHJ7fb3N6LaO3vm/3V5o+YMlyeZ0gNbEIMracUDLZvYQUNnMYchd+kfeM2zfVVBts7uCKjjPaNsHX3yIpjSAWyd+0S0g3ptnqUwjI9PMcP6W1tkWzdQa3NEyzOU0QaOFLspeNDRPZq2FlCv9uhFs4R+C5pUpiwGaeF8e04KUxZXMOi69q4chgms3TIDpM8NyVqncJqEC1G1/qmtq/5RMFAb1QDqyWMtge9tgHYQdQnzurIRJtUQWFVCAKPtMiJdvfI6xLbDyhzgRcEYHmUUhi6r8OfKBMD+vJ2t0vDkabErRethQHuIDG+kyUZRWmxf+suYS1ge+sOJ06uIqzYPNTJYtOf7x6k2I5liJDeaKzg4CA2wOZqaSOMxkv0gTKyLKdIM6TjUDq2sR7N5zWD1WmtDm9mvFLzCcclyYba1owvzYd7073BfhShipxaxTcjr4XwTRjXzZBqxYF4QOjk5CJhkJUErksZdSnihKwsjGXJj268wflghXxgsby4hFQFUjkG2Ow8JS8lt+7eJb91G7s+h//Bbcp2C9d5lb5VZ+/2HayFKq7rMtDNU8cxPimdof/per5GbW0B+weZcQUjqEASDxKkI5G+awqXdd9mkChCG5bWT5hkR7vTXD2g14tpdzvGKsateI30QjoM4gTsEIJ54oMWXpGS7N0hO+hyN+4T9VOqFZ+tbs/wh9buPrrIsdtXyF4y4Gp8B7lfmjy/2pjH83VWaJtqb5YW+EGVb31wBVdc5ScePsJOu0fNEeSLVYq4h0q1CRdE3YHhEfohGqGlPqhOjUYMUaO+6RaPZg+qVR/PtYgGqalJKHc4/BAnJZXl06xufIYz1Q413wLpsnV3j5ffesvgiLGGNKFaqZAOEnb2D8iziIWGz4c/3GL9xDE+fO+aKYtt7seEns2g06WfDGsPOmKdPrmAdKKcrajgSCOgXexBXOcgEaRJxs5BTnMxI4g6PPf4Y3z9e9/n977xEUtLDTzpUnEyqoHUXQ/TDarVQqJBbHoD1XoVW9qaxRL1E7IkJdAlKdeejK5pK9GuUAmlaabaFNQr0pi4lWQsr5yEcofV9QoVp8D3Krz69jsGMHUYaXf2jaWuNuvcurvDx7EOezGXbu6yPH+JbqSG1F5HmLw0OFdvhCQ6LLsW71+7i9T03gltWknBimMhsi0q8/fjhBZhEZseu1upGn/9yc9/gRde+g5X220K22Ej6BD6RxjoXCCOyXQUiRJc3ycexOT6cPXAJEoycE346x0MDBnS2KATHC0QHXl0bB/YAllkQyDNPSpr55m3V3AbYKs+jaaN1BEnG0zmVrrdDhXXoRPn5Fjc2hzgVwJu9UIqoTucHYp67B8ozp+o8eGmYP34Ct1BBtEBspoPJV/zbZZcgRX3yHbep5x7gK42TR2+lCJKOhSqx/NPP8HrH33Ize3rxIlCaD9Uka6sU6iSoB4Mc381HNAZ6Bqd1Jmii62txAlMn0Ff0+ivGaR2t9DRxRAxjA5piWWDIzKk5xMXKf2iSbu0qcwdId6+PJkVNuVRUXDl7haubQb/iHOoegW2Gk6K9cuCICy4cSDwZAmZ3keBrA2QiWsxnxc0o5h+J8auBgRlSiX+kNA7hipSQ1acSo0oKejHCY/d/yCBX6G1d5msSMk0MluusRI/dI3/65mlNFU40jdZYeCPhps0AVIlgzgnzpUhM5o9akDT4XNI2S0cX+LriORCik9eJBRC4FcXEXevTOaWdBuvFnrGbaRvUW84dPYTlud9k4xJ3anWUU24pmWmLUIVHYSQ+L6NvJ2AjHPD42u6YNEf0CsXjGbnnAOOLiR0+x6bt/eYa9ZM/qzT3wunznGtNket+wC+t0SqabRtmSRFaSDUh5DmLXNAzQkMf7CH3+jQ+BBqXqCG31LR4Knf11akQ5v+TVClV7pYKjUgq6RLY/kY7buXyJKB4Sxm5tC26UY6DbWIuwWDgSBvJ6NrgkKbhF2S6+aqLtRYFrYcTq7IM7pxWSpe3YQy02zM4rkvnefM/B3OnHsKaie5cbmDvzcwwOgHAUkckauE1UqNr84P+GY7o6fRzozVKTyN2rriUwy/g2DyczMap0xLW1P02CpNUqKrwroAbUpwOj+vSEOQBv0CkSX4dkrg5eYzuusbNkKqi8fI+rvD9ntZkMsK/TIg7iqkLEmTkkxzfb2wDVmsm6nDL11pF/MCicgFWVIgz596yPD7c2djQgsGBdSay9yMbKxP+vTVZTZvtdnaPsB2HGqVlPl6lTCwCQOPpurw80cV//v2Eh+3hzm31rT2ee0GOhGyxWh8zQI7G46t6xxB39cbDKs92i2USo0LSFe31R2sikMmCqJcoTQgN33i65cpVYQThgY7zBSa7xHWfGN5uscofAdsj9Ie5R1eYQQhyRmkFsKRxhriIkdWTj5hcvU528Kp1XHCCqX26/6ANIqo2BbzjQXuP5OZDRversvnYYjray4+1PwvrnX5D6+tUJTCfLdAd431bNA4ZtumICEojNkyrBhZlunUaHzw5BAAx7NGBSX9foxTxvh2ZmqNbU1x8xRbJzVJD+FK43Je6LA4VzdYo3sPRZ4TpTlRLkySpa1Qy0qzTD8cMlaTtRYK+dHmVfppbAjDajUgdENjOt0iMZmhKZvrWbwswtMqzFJy4TGvFs37GYEAAAOaSURBVPHSCrbMyZIcRUaz9Im0xm2HVi4plW00a+K+RnUp0DTIfFFLaa0OEVyHBKMZExpGJTd7KCTH8gk0sSoiMk+RpxF5dDAqsAzvHTdJA9eiWama39ri9k1RV5CpkqonkTo71f0ePXkiShaqDvIvrr5CotvgOi0udFwuDcgNHIHStLvMsUo99aUTIBtfCO5EKRcbPj92+inm62fxK67ps/3iI+N2Fxz0Ez7Zt3hp02crkUNfzUuyQtFPc/b3tmnvXKUoUjNg4XoV6rVlgkqNSqVuSJL2mUxYxozTtKCVuIZ/CC8YRg7XMSSnW1i0ezknl2smvC4F0iRg1WBYqtMzSqGZJslNl0rn5Hpqbb+XYv+tr5z+mpK2pcOF1CRcC7/qUs0sajrGOy6Z0JkddGVJ6FbZqC+xmXX5wfZ19pM9lv0GgeeboURT8pICX8JqCA81EjTDvRX7phylOb/OCoOwSrWxQrU2b8bpbHIcKzNFlHpDD23lRkueU+ruDUKlJIOe6RrrabJhrjQcga7Ikucu6C4zpowXJ6lppO52IvpxRi9KTKVL5yQ6HdbVKj3ItVhzlUyzhJVCsJBDKxfcqfhs9xTKyUlFaUwYmeEsCWquw3anxa3BDoGrwLf5ML5C9+4mX1JfphYuYGv+L6QZaJI69/ZdPu/HXKz1+fZunQ87rmmI2mbAeVhP9MImflAnbd+g0z8gOmgTRftIv2msoV71aG3fodvdpcz7lHk2TIFdx/QewsAxPj9X1emwMoNUCcL0NSuexDFfsRmW5sffLrFGIdc+9Wz9a7fJrLueYoeETpEwSLtkMjaUNE1zMwpbiMKgrJ7kzK3hlxh04SIflLSihCvxxxS9NnY/wNbdGs82GrQdG7/i06i6nAliAlLe3y1IlYXnuzTmGlQ0oLoublDR26KwbGyvaipT2nWuffQqUWeLUvVN604GDoXnUUjJWt3jcxs1GoFtSJXmKaFjEzhDTPEcaTTu2MPJdX1oHa37ulaRKSUX604eq8z4hwrMNygIM8FBJlC2QlfqvUAw73g0qh4HusCoi6e5wNHhzZZkomRrd8DrzSv41hy1Ym44vGAPv0FuOkM6+fEcHluMWbATfuudO7x77TYqj5GeZ6zAllUzKJ0mXdPRLUViDrBx/8NYXoXB/i6D1g6d7hVK7cu5YrNT8uIHBT/zaJOGrnq5tinFaUFo0qVHfPTZB2lhmiwaBPV+NN9xhMj/Pyp9uiMrMrjfAAAAAElFTkSuQmCC"

	// Create a reader to split the encoding.
	reader := strings.NewReader(imgBase64)

	// Split the encoding into multiple manage data operations.
	var ops []txnbuild.Operation
	index := 0
	name := make([]byte, 62)
	value := make([]byte, 64)
	for reader.Len() > 0 {
		// Trim slices to amount of bytes read.
		count, _ := reader.Read(name)
		name = name[:count]
		count, _ = reader.Read(value)
		value = value[:count]

		// Append the operation.
		ops = append(ops, &txnbuild.ManageData{
			Name: fmt.Sprintf("%02d", index) + string(name),
			// Not double casting the value byte array causes the
			// transmission to send the same value string for each key.
			Value: []byte(string(value)),
		})
		index++
	}

	// Note that the transaction is split to 2, otherwise
	// we get "HTTP 414 URI Too Long" error. Check issue
	// #3621 on stellar/go for more info.

	// Construct the transaction (1/2).
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops[:len(ops)/2],
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Sign the transaction (1/2).
	tx, err = tx.Sign(network.TestNetworkPassphrase, questAccount.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	// Send the transaction to the network (1/2).
	status, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the response.
	fmt.Printf("Successfully submitted transaction!\nTransaction ID: %v\n", status.ID)

	// Construct the transaction (2/2).
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops[len(ops)/2:],
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Sign the transaction (2/2).
	tx, err = tx.Sign(network.TestNetworkPassphrase, questAccount.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	// Send the transaction to the network (2/2).
	status, err = client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the response.
	fmt.Printf("Successfully submitted transaction!\nTransaction ID: %v\n", status.ID)

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
