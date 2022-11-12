package dbrequests

import "go.mongodb.org/mongo-driver/bson"

// AddRate generate mongodb request for adding rate
func AddRate(rateName string) bson.A {
	return bson.A{
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "Rating." + rateName,
						Value: bson.D{
							{Key: "$sum",
								Value: bson.A{
									"$Rating." + rateName,
									1,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "Rating.Average",
						Value: bson.D{
							{Key: "$divide",
								Value: bson.A{
									bson.D{
										{Key: "$sum",
											Value: bson.A{
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.Five",
															5,
														},
													},
												},
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.Four",
															4,
														},
													},
												},
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.Three",
															3,
														},
													},
												},
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.Two",
															2,
														},
													},
												},
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.One",
															1,
														},
													},
												},
											},
										},
									},
									bson.D{
										{Key: "$sum",
											Value: bson.A{
												"$Rating.Five",
												"$Rating.Four",
												"$Rating.Three",
												"$Rating.Two",
												"$Rating.One",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// ChangeRating generate mongodb request for changing rate
func ChangeRating(from string, to string) bson.A {
	return bson.A{
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "Rating." + from,
						Value: bson.D{
							{Key: "$sum",
								Value: bson.A{
									"$Rating." + from,
									-1,
								},
							},
						},
					},
					{Key: "Rating." + to,
						Value: bson.D{
							{Key: "$sum",
								Value: bson.A{
									"$Rating." + to,
									1,
								},
							},
						},
					},
				},
			},
		},
		bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "Rating.Average",
						Value: bson.D{
							{Key: "$divide",
								Value: bson.A{
									bson.D{
										{Key: "$sum",
											Value: bson.A{
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.Five",
															5,
														},
													},
												},
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.Four",
															4,
														},
													},
												},
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.Three",
															3,
														},
													},
												},
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.Two",
															2,
														},
													},
												},
												bson.D{
													{Key: "$multiply",
														Value: bson.A{
															"$Rating.One",
															1,
														},
													},
												},
											},
										},
									},
									bson.D{
										{Key: "$sum",
											Value: bson.A{
												"$Rating.Five",
												"$Rating.Four",
												"$Rating.Three",
												"$Rating.Two",
												"$Rating.One",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
